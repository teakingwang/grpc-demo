package service

import (
	"context"
	"database/sql"

	"github.com/teakingwang/grpc-demo/pkg/errors"
	pb "github.com/teakingwang/grpc-demo/proto/user/gen"

	"golang.org/x/crypto/bcrypt"
)

// UserService 实现用户服务接口
type UserService struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB // 数据库连接
}

// NewUserService 创建用户服务实例
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// 对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 将用户信息插入数据库
	result, err := s.db.ExecContext(ctx,
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		req.Username, req.Email, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	// 获取插入的用户ID
	id, _ := result.LastInsertId()
	return &pb.CreateUserResponse{
		Id:       id,
		Username: req.Username,
	}, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user pb.GetUserResponse

	// 从数据库查询用户信息
	err := s.db.QueryRowContext(ctx,
		"SELECT id, username, email FROM users WHERE id = ?",
		req.Id).Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.FromError(err)
	}
	return &user, nil
}
