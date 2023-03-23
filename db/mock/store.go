// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/asdsec/thenut/db/sqlc (interfaces: Store)

// Package mock_db is a generated GoMock package.
package mock_db

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	db "github.com/asdsec/thenut/db/sqlc"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddMerchantBalance mocks base method.
func (m *MockStore) AddMerchantBalance(arg0 context.Context, arg1 db.AddMerchantBalanceParams) (db.Merchant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMerchantBalance", arg0, arg1)
	ret0, _ := ret[0].(db.Merchant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMerchantBalance indicates an expected call of AddMerchantBalance.
func (mr *MockStoreMockRecorder) AddMerchantBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMerchantBalance", reflect.TypeOf((*MockStore)(nil).AddMerchantBalance), arg0, arg1)
}

// CreateConsultancy mocks base method.
func (m *MockStore) CreateConsultancy(arg0 context.Context, arg1 db.CreateConsultancyParams) (db.Consultancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateConsultancy", arg0, arg1)
	ret0, _ := ret[0].(db.Consultancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateConsultancy indicates an expected call of CreateConsultancy.
func (mr *MockStoreMockRecorder) CreateConsultancy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateConsultancy", reflect.TypeOf((*MockStore)(nil).CreateConsultancy), arg0, arg1)
}

// CreateCustomer mocks base method.
func (m *MockStore) CreateCustomer(arg0 context.Context, arg1 string) (db.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCustomer", arg0, arg1)
	ret0, _ := ret[0].(db.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCustomer indicates an expected call of CreateCustomer.
func (mr *MockStoreMockRecorder) CreateCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomer", reflect.TypeOf((*MockStore)(nil).CreateCustomer), arg0, arg1)
}

// CreateEntry mocks base method.
func (m *MockStore) CreateEntry(arg0 context.Context, arg1 db.CreateEntryParams) (db.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEntry", arg0, arg1)
	ret0, _ := ret[0].(db.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEntry indicates an expected call of CreateEntry.
func (mr *MockStoreMockRecorder) CreateEntry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEntry", reflect.TypeOf((*MockStore)(nil).CreateEntry), arg0, arg1)
}

// CreateMerchant mocks base method.
func (m *MockStore) CreateMerchant(arg0 context.Context, arg1 db.CreateMerchantParams) (db.Merchant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMerchant", arg0, arg1)
	ret0, _ := ret[0].(db.Merchant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMerchant indicates an expected call of CreateMerchant.
func (mr *MockStoreMockRecorder) CreateMerchant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMerchant", reflect.TypeOf((*MockStore)(nil).CreateMerchant), arg0, arg1)
}

// CreatePayment mocks base method.
func (m *MockStore) CreatePayment(arg0 context.Context, arg1 db.CreatePaymentParams) (db.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePayment", arg0, arg1)
	ret0, _ := ret[0].(db.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePayment indicates an expected call of CreatePayment.
func (mr *MockStoreMockRecorder) CreatePayment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePayment", reflect.TypeOf((*MockStore)(nil).CreatePayment), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// DeleteCustomer mocks base method.
func (m *MockStore) DeleteCustomer(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCustomer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCustomer indicates an expected call of DeleteCustomer.
func (mr *MockStoreMockRecorder) DeleteCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCustomer", reflect.TypeOf((*MockStore)(nil).DeleteCustomer), arg0, arg1)
}

// DeleteMerchant mocks base method.
func (m *MockStore) DeleteMerchant(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMerchant", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMerchant indicates an expected call of DeleteMerchant.
func (mr *MockStoreMockRecorder) DeleteMerchant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMerchant", reflect.TypeOf((*MockStore)(nil).DeleteMerchant), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1)
}

// GetConsultancy mocks base method.
func (m *MockStore) GetConsultancy(arg0 context.Context, arg1 int64) (db.Consultancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConsultancy", arg0, arg1)
	ret0, _ := ret[0].(db.Consultancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConsultancy indicates an expected call of GetConsultancy.
func (mr *MockStoreMockRecorder) GetConsultancy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConsultancy", reflect.TypeOf((*MockStore)(nil).GetConsultancy), arg0, arg1)
}

// GetCustomer mocks base method.
func (m *MockStore) GetCustomer(arg0 context.Context, arg1 int64) (db.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomer", arg0, arg1)
	ret0, _ := ret[0].(db.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCustomer indicates an expected call of GetCustomer.
func (mr *MockStoreMockRecorder) GetCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomer", reflect.TypeOf((*MockStore)(nil).GetCustomer), arg0, arg1)
}

// GetEntry mocks base method.
func (m *MockStore) GetEntry(arg0 context.Context, arg1 int64) (db.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntry", arg0, arg1)
	ret0, _ := ret[0].(db.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntry indicates an expected call of GetEntry.
func (mr *MockStoreMockRecorder) GetEntry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntry", reflect.TypeOf((*MockStore)(nil).GetEntry), arg0, arg1)
}

// GetMerchant mocks base method.
func (m *MockStore) GetMerchant(arg0 context.Context, arg1 int64) (db.Merchant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMerchant", arg0, arg1)
	ret0, _ := ret[0].(db.Merchant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMerchant indicates an expected call of GetMerchant.
func (mr *MockStoreMockRecorder) GetMerchant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMerchant", reflect.TypeOf((*MockStore)(nil).GetMerchant), arg0, arg1)
}

// GetPayment mocks base method.
func (m *MockStore) GetPayment(arg0 context.Context, arg1 int64) (db.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPayment", arg0, arg1)
	ret0, _ := ret[0].(db.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPayment indicates an expected call of GetPayment.
func (mr *MockStoreMockRecorder) GetPayment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPayment", reflect.TypeOf((*MockStore)(nil).GetPayment), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// ListConsultancies mocks base method.
func (m *MockStore) ListConsultancies(arg0 context.Context, arg1 db.ListConsultanciesParams) ([]db.Consultancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListConsultancies", arg0, arg1)
	ret0, _ := ret[0].([]db.Consultancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListConsultancies indicates an expected call of ListConsultancies.
func (mr *MockStoreMockRecorder) ListConsultancies(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListConsultancies", reflect.TypeOf((*MockStore)(nil).ListConsultancies), arg0, arg1)
}

// ListEntries mocks base method.
func (m *MockStore) ListEntries(arg0 context.Context, arg1 db.ListEntriesParams) ([]db.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEntries", arg0, arg1)
	ret0, _ := ret[0].([]db.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEntries indicates an expected call of ListEntries.
func (mr *MockStoreMockRecorder) ListEntries(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntries", reflect.TypeOf((*MockStore)(nil).ListEntries), arg0, arg1)
}

// ListMerchants mocks base method.
func (m *MockStore) ListMerchants(arg0 context.Context, arg1 db.ListMerchantsParams) ([]db.Merchant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMerchants", arg0, arg1)
	ret0, _ := ret[0].([]db.Merchant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMerchants indicates an expected call of ListMerchants.
func (mr *MockStoreMockRecorder) ListMerchants(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMerchants", reflect.TypeOf((*MockStore)(nil).ListMerchants), arg0, arg1)
}

// ListPayments mocks base method.
func (m *MockStore) ListPayments(arg0 context.Context, arg1 db.ListPaymentsParams) ([]db.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPayments", arg0, arg1)
	ret0, _ := ret[0].([]db.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPayments indicates an expected call of ListPayments.
func (mr *MockStoreMockRecorder) ListPayments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPayments", reflect.TypeOf((*MockStore)(nil).ListPayments), arg0, arg1)
}

// UpdateCustomer mocks base method.
func (m *MockStore) UpdateCustomer(arg0 context.Context, arg1 db.UpdateCustomerParams) (db.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCustomer", arg0, arg1)
	ret0, _ := ret[0].(db.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCustomer indicates an expected call of UpdateCustomer.
func (mr *MockStoreMockRecorder) UpdateCustomer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCustomer", reflect.TypeOf((*MockStore)(nil).UpdateCustomer), arg0, arg1)
}

// UpdateEmail mocks base method.
func (m *MockStore) UpdateEmail(arg0 context.Context, arg1 db.UpdateEmailParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEmail indicates an expected call of UpdateEmail.
func (mr *MockStoreMockRecorder) UpdateEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockStore)(nil).UpdateEmail), arg0, arg1)
}

// UpdateMerchant mocks base method.
func (m *MockStore) UpdateMerchant(arg0 context.Context, arg1 db.UpdateMerchantParams) (db.Merchant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMerchant", arg0, arg1)
	ret0, _ := ret[0].(db.Merchant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMerchant indicates an expected call of UpdateMerchant.
func (mr *MockStoreMockRecorder) UpdateMerchant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMerchant", reflect.TypeOf((*MockStore)(nil).UpdateMerchant), arg0, arg1)
}

// UpdatePassword mocks base method.
func (m *MockStore) UpdatePassword(arg0 context.Context, arg1 db.UpdatePasswordParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockStoreMockRecorder) UpdatePassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockStore)(nil).UpdatePassword), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockStore) UpdateUser(arg0 context.Context, arg1 db.UpdateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStoreMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStore)(nil).UpdateUser), arg0, arg1)
}
