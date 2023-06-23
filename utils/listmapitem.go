package utils

import "kubmanager/model/base"

func ToMap(items []base.ListMapItem) map[string]string {
	dataMap := make(map[string]string)
	for _, item := range items {
		dataMap[item.Key] = item.Value
	}
	return dataMap
}

func ToList(item map[string]string) []base.ListMapItem {
	itemList := make([]base.ListMapItem, 0)
	for k, v := range item {
		itemList = append(itemList, base.ListMapItem{
			Key:   k,
			Value: v,
		})

	}
	return itemList
}

func ToListWithMapByte(data map[string][]byte) []base.ListMapItem {
	itemList := make([]base.ListMapItem, 0)
	for k, v := range data {
		itemList = append(itemList, base.ListMapItem{
			Key:   k,
			Value: string(v),
		})

	}
	return itemList
}
