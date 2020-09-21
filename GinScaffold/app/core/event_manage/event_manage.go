package event_manage

import (
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"sync"
	"strings"
)

// 定义一个全局事件存储变量，本模块只负责存储 键 => 函数 ， 相对容器来说功能稍弱，但是调用更加简单、方便、快捷
// sync.Map 无需初始化
// http://c.biancheng.net/view/34.html
var sMap sync.Map

// 创建一个事件管理工厂：event_manage都是使用的时候创建的，但是操作的都是 sMap 这个全局变量
func CreateEventManageFactory() *eventManage {
	return &eventManage{}
}

type eventManage struct {
}

// 1. 注册时间
func (e *eventManage) Set(key string, keyFunc func(args ...interface{})) bool {
	if _, exists := e.Get(key); exists == false {
		sMap.Store(key, keyFunc)
		return true
	} else {
		variable.ZapLog.Info(my_errors.ErrorsFuncEventAlreadyExists + " ,相关键名：" + key)
	}
	return false
}

// 2. 获取时间
func (e *eventManage) Get(key string) (interface{}, bool) {
	if value, exists := sMap.Load(key); exists {
		return value, exists
	}
	return nil, false
}

// 3. 执行事件
func (e *eventManage) Call(key string, args ...interface{}) {
	if valueInterface, exists := sMap.Load(key); exists {
		if fn, ok := valueInterface.(func(args ...interface{})); ok {
			fn(args...)
		} else {
			variable.ZapLog.Error(my_errors.ErrorsFuncEventNotCall + " ,键名：" + key + ", 相关函数无法调用")
		}
	} else {
		variable.ZapLog.Error(my_errors.ErrorsFuncEventNotRegister + " ,键名：" + key)
	}
}

// 4. 删除事件
func (e *eventManage) Delete(key string) {
	sMap.Delete(key)
}

// 5. 根据键的前缀，模糊调用，请慎重使用
func (e *eventManage) FuzzyCall(keyPre string) {
	sMap.Range(func(key, value interface{}) bool {
		if keyName, ok := key.(string); ok {
			if strings.HasPrefix(keyName, keyPre) {
				e.Call(keyName)
			}
		}
		return true
	})
}



