package backend

import (
	"reflect"
	"testing"

	"facette.io/sliceutil"
	"facette.io/sqlstorage"
	"github.com/stretchr/testify/assert"
)

func testInterfaceToSlice(v interface{}) []interface{} {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		return nil
	}

	result := []interface{}{}
	for i, n := 0, rv.Len(); i < n; i++ {
		result = append(result, rv.Index(i).Interface())
	}

	return result
}

func testNewItem(v interface{}) interface{} {
	return reflect.New(reflect.TypeOf(v).Elem()).Interface()
}

func testNewItemSlice(v interface{}) interface{} {
	return reflect.New(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(v)), 0, 0).Type()).Interface()
}

func testItemCreate(b *Backend, refItem interface{}, testItems []interface{}, t *testing.T) {
	assert.Nil(t, b.Storage().Save(testItems[0]))
}

func testItemCreateInvalid(b *Backend, refItem interface{}, testItems []interface{}, t *testing.T) {
	item := testNewItem(refItem)
	reflect.Indirect(reflect.ValueOf(item)).FieldByName("Name").SetString("invalid!")
	assert.Equal(t, ErrInvalidName, b.Storage().Save(item))

	item = testNewItem(refItem)
	reflect.Indirect(reflect.ValueOf(item)).FieldByName("ID").SetString("invalid!")
	assert.Equal(t, ErrInvalidID, b.Storage().Save(item))
}

func testItemGet(b *Backend, refItem interface{}, testItems []interface{}, t *testing.T) {
	item := testNewItem(refItem)
	rv := reflect.Indirect(reflect.ValueOf(item))
	assert.Nil(t, b.Storage().Get("name", "item1", item, true))
	assert.NotZero(t, rv.FieldByName("ID").String())
	assert.NotZero(t, rv.FieldByName("Created").Interface())
	assert.Equal(t, testItems[0], item)
}

func testItemGetUnknown(b *Backend, refItem interface{}, testItems []interface{}, t *testing.T) {
	item := testNewItem(refItem)
	assert.Equal(t, sqlstorage.ErrItemNotFound, b.Storage().Get("name", "unknown1", item, true))
}

func testItemUpdate(b *Backend, refItem interface{}, testItems []interface{}, t *testing.T) {
	desc := "A great item description"

	dv := reflect.Indirect(reflect.ValueOf(testItems[0])).FieldByName("Description")
	dv.Set(reflect.ValueOf(&desc))

	assert.Nil(t, b.Storage().Save(testItems[0]))

	item := testNewItem(refItem)
	rv := reflect.Indirect(reflect.ValueOf(item))
	assert.Nil(t, b.Storage().Get("name", "item1", item, true))
	assert.NotZero(t, rv.FieldByName("ID").String())
	assert.NotZero(t, rv.FieldByName("Created").Interface())
	assert.NotZero(t, rv.FieldByName("Modified").Interface())
	assert.Equal(t, testItems[0], item)

	desc = ""
	dv.Set(reflect.ValueOf(&desc))

	assert.Nil(t, b.Storage().Save(testItems[0]))

	item = testNewItem(refItem)
	rv = reflect.Indirect(reflect.ValueOf(item))
	assert.Nil(t, b.Storage().Get("name", "item1", item, true))
	assert.Nil(t, rv.FieldByName("Description").Interface())
	assert.Equal(t, testItems[0], item)
}

func testItemCount(b *Backend, refItem interface{}, testItems []interface{}, t *testing.T) {
	count, err := b.Storage().Count(testNewItem(refItem))
	assert.Nil(t, err)
	assert.Equal(t, 3, count)
}

func testItemList(b *Backend, refItem interface{}, testItems []interface{}, t *testing.T) {
	for _, item := range testItems {
		assert.Nil(t, b.Storage().Save(item))
	}

	items := testNewItemSlice(refItem)
	count, err := b.Storage().List(items, nil, []string{"name"}, 0, 0, true)
	assert.Nil(t, err)
	assert.Equal(t, len(testItems), count)
	assert.Equal(t, testItems, testInterfaceToSlice(items))

	items = testNewItemSlice(refItem)
	count, err = b.Storage().List(items, map[string]interface{}{"name": "item1"}, []string{"name"}, 0, 0, true)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
	assert.Equal(t, testItems[:1], testInterfaceToSlice(items))

	items = testNewItemSlice(refItem)
	count, err = b.Storage().List(items, nil, []string{"name"}, 0, 2, true)
	assert.Nil(t, err)
	assert.Equal(t, len(testItems), count)
	assert.Equal(t, testItems[:2], testInterfaceToSlice(items))

	items = testNewItemSlice(refItem)
	count, err = b.Storage().List(items, nil, []string{"name"}, 1, 1, true)
	assert.Nil(t, err)
	assert.Equal(t, len(testItems), count)
	assert.Equal(t, testItems[1:2], testInterfaceToSlice(items))

	sliceutil.Reverse(testItems)

	items = testNewItemSlice(refItem)
	count, err = b.Storage().List(items, nil, []string{"-name"}, 0, 0, true)
	assert.Nil(t, err)
	assert.Equal(t, len(testItems), count)
	assert.Equal(t, testItems, testInterfaceToSlice(items))

	items = testNewItemSlice(refItem)
	count, err = b.Storage().List(items, nil, []string{"-name"}, 2, 10, true)
	assert.Nil(t, err)
	assert.Equal(t, len(testItems), count)
	assert.Equal(t, testItems[2:], testInterfaceToSlice(items))
}

func testItemDelete(b *Backend, refItem interface{}, testItems []interface{}, t *testing.T) {
	assert.Nil(t, b.Storage().Delete(testItems[0]))
	assert.Equal(t, sqlstorage.ErrItemNotFound, b.Storage().Get("name", "item1", testNewItem(refItem), true))
}

func testItemDeleteAll(b *Backend, refItem interface{}, testItems []interface{}, t *testing.T) {
	assert.Nil(t, b.Storage().Delete(testNewItem(refItem)))
}
