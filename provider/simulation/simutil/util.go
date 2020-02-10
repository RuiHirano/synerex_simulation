package simutil

import (
	"log"
	"sync"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(format string, args ...interface{}) {
	log.Printf("\x1b[32m\x1b[40m \n [Info] \x1b[0m"+format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	log.Printf("\x1b[31m\x1b[40m \n [Error] \x1b[0m"+format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	log.Printf("\x1b[33m\x1b[40m \n [Warn] \x1b[0m"+format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	log.Fatalf("\x1b[31m\x1b[40m \n [Error] \x1b[0m"+format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	log.Printf("\x1b[36m\x1b[40m \n [Debug] \x1b[0m"+format, args...)
}

type SyncMap struct {
	mapData *sync.Map
}

func NewMap() *SyncMap {
	return &SyncMap{mapData: new(sync.Map)}
}

// StoreMap mapにデータをストアする
func (thisMap *SyncMap) Store(key, value interface{}) {
	thisMap.mapData.Store(key, value)
}

// LoadMap mapからデータを取得する
func (thisMap *SyncMap) Load(key interface{}) (value interface{}, ok bool) {
	return thisMap.mapData.Load(key)
}

// Delete mapにあるデータを削除
func (thisMap *SyncMap) Delete(key interface{}) {
	thisMap.mapData.Delete(key)
}
