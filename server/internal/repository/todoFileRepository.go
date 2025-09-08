package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"hw25/internal/model"
	"io"
	"os"
	"strings"
)

func readFromFile[T model.HomeworkItem | model.StudyItem | model.WorkoutItem](fileName string) ([]*T, error) {
	var result []*T
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}

		return nil, errors.New(fmt.Sprintf("Cannot open file %s, err: %s", fileName, err.Error()))
	}
	defer file.Close()

	writer := new(strings.Builder)
	io.Copy(writer, file)
	reader := strings.NewReader(writer.String())

	decoder := json.NewDecoder(reader)

	for {
		var item T
		err := decoder.Decode(&item)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Cannot decode item from file %s, err: %s", fileName, err.Error()))
		}
		result = append(result, &item)
	}

	return result, nil
}

func appendToFile[T *model.HomeworkItem | *model.StudyItem | *model.WorkoutItem](fileName string, dataToAppend T) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot open file %s, err: %s", fileName, err.Error()))
	}
	defer file.Close()

	buf, err := json.Marshal(dataToAppend)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot marshal data in order to write to file %s, err: %s", fileName, err.Error()))
	}

	_, err = file.Write(buf)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot write data to file %s, err: %s", fileName, err.Error()))
	}

	return nil
}

func overwriteFile[T *model.HomeworkItem | *model.StudyItem | *model.WorkoutItem](fileName string, slice []T) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot open file %s, err: %s", fileName, err.Error()))
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 64*1024)
	defer writer.Flush()

	for _, item := range slice {
		data, err := json.Marshal(item)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot marshal data in order to write to file %s, err: %s", fileName, err.Error()))
		}

		if _, err := writer.Write(data); err != nil {
			return errors.New(fmt.Sprintf("Cannot write data to file %s, err: %s", fileName, err.Error()))
		}
	}

	return nil
}
