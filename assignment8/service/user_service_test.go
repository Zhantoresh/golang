package service

import (
	"errors"
	"assignment8/repository"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockDB)

	expectedUsr := &repository.User{ID: 10, Name: "Test Client"}
	mockDB.EXPECT().GetUserByID(10).Return(expectedUsr, nil)

	actual, err := svc.GetUserByID(10)
	require.NoError(t, err)
	require.Equal(t, expectedUsr, actual)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockDB)

	newUsr := &repository.User{ID: 5, Name: "Newbie"}
	mockDB.EXPECT().CreateUser(newUsr).Return(nil)

	err := svc.CreateUser(newUsr)
	require.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockDB)
	
	targetEmail := "admin@domain.com"
	usr := &repository.User{ID: 7, Name: "Admin"}

	t.Run("user_exists_conflict", func(t *testing.T) {
		mockDB.EXPECT().GetByEmail(targetEmail).Return(usr, nil)
		err := svc.RegisterUser(usr, targetEmail)
		require.ErrorContains(t, err, "already exists")
	})

	t.Run("successful_registration", func(t *testing.T) {
		mockDB.EXPECT().GetByEmail(targetEmail).Return(nil, nil)
		mockDB.EXPECT().CreateUser(usr).Return(nil)
		err := svc.RegisterUser(usr, targetEmail)
		require.NoError(t, err)
	})

	t.Run("db_failure_on_creation", func(t *testing.T) {
		mockDB.EXPECT().GetByEmail(targetEmail).Return(nil, nil)
		mockDB.EXPECT().CreateUser(usr).Return(errors.New("connection lost"))
		err := svc.RegisterUser(usr, targetEmail)
		require.EqualError(t, err, "connection lost")
	})
}

func TestUpdateUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockDB)

	targetID := 42
	existingUsr := &repository.User{ID: targetID, Name: "OldName"}

	t.Run("validation_empty_name", func(t *testing.T) {
		err := svc.UpdateUserName(targetID, "")
		require.EqualError(t, err, "name cannot be empty")
	})

	t.Run("user_not_found_in_db", func(t *testing.T) {
		mockDB.EXPECT().GetUserByID(targetID).Return(nil, errors.New("not found"))
		err := svc.UpdateUserName(targetID, "NewName")
		require.Error(t, err)
	})

	t.Run("successful_name_update", func(t *testing.T) {
		mockDB.EXPECT().GetUserByID(targetID).Return(existingUsr, nil)
		mockDB.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(u *repository.User) error {
			require.Equal(t, "FreshName", u.Name)
			return nil
		})
		
		err := svc.UpdateUserName(targetID, "FreshName")
		require.NoError(t, err)
	})

	t.Run("update_fails_in_repo", func(t *testing.T) {
		existingUsr.Name = "OldName" 
		mockDB.EXPECT().GetUserByID(targetID).Return(existingUsr, nil)
		mockDB.EXPECT().UpdateUser(existingUsr).Return(errors.New("tx failed"))
		err := svc.UpdateUserName(targetID, "FreshName")
		require.EqualError(t, err, "tx failed")
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockDB)

	t.Run("protect_admin_deletion", func(t *testing.T) {
		err := svc.DeleteUser(1)
		require.ErrorContains(t, err, "not allowed to delete admin user")
	})

	t.Run("delete_success", func(t *testing.T) {
		mockDB.EXPECT().DeleteUser(99).Return(nil)
		err := svc.DeleteUser(99)
		require.NoError(t, err)
	})

	t.Run("repo_throws_error", func(t *testing.T) {
		mockDB.EXPECT().DeleteUser(99).Return(errors.New("db locked"))
		err := svc.DeleteUser(99)
		require.EqualError(t, err, "db locked")
	})
}