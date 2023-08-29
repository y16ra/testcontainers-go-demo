package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/y16ra/testcontainers-go-demo/demo/model"
)

func Test_userRepository_FindById(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				id: 1,
			},
			want: &model.User{
				ID:   1,
				Name: "test",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	db, err := setupDBContainer(t)
	initDataForFindById(db)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{
				db: db,
			}
			got, err := r.FindById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(model.User{}),
				cmpopts.IgnoreFields(model.User{}, "CreatedAt"),
			}
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("missmatch (got, want) :%s", diff)
			}

		})
	}
}

func initDataForFindById(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (id, name) VALUES ($1, $2)", 1, "test")
	if err != nil {
		return err
	}
	return nil
}

func Test_userRepository_Store(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				user: &model.User{
					Name:      "test",
					CreatedAt: time.Now(),
				},
			},
			wantErr: false,
		},
	}
	db, err := setupDBContainer(t)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{
				db: db,
			}
			if err := r.Store(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("userRepository.Store() user = %v", tt.args.user)
		})
	}
}
