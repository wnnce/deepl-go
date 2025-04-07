package deepl

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewCMD(t *testing.T) {
	cmd := NewCMD(context.Background(), func() (int, error) {
		return 0, nil
	})
	fmt.Println(cmd.Closed())
}

func TestCMD_Sync(t *testing.T) {
	cmd := NewCMD(context.Background(), func() (int, error) {
		time.Sleep(3 * time.Second)
		return 1, nil
	})
	result, err := cmd.Sync()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(cmd.Closed())
	fmt.Println(result)
}

func TestCMD_Async(t *testing.T) {
	cmd := NewCMD(context.Background(), func() (int, error) {
		time.Sleep(3 * time.Second)
		return 100, nil
	})
	cmd.Async(func(ctx context.Context, result int, err error) {
		fmt.Println(result)
		fmt.Println(cmd.Closed())
	})
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		time.Sleep(500 * time.Millisecond)
	}
}
