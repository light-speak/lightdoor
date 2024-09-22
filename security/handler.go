package main

import (
	"context"

	"github.com/light-speak/lightdoor/security/jwt"
	token "github.com/light-speak/lightdoor/security/kitex_gen/token"
)

// SecurityServiceImpl implements the last service interface defined in the IDL.
type SecurityServiceImpl struct{}

// GetSecurityUserId implements the SecurityServiceImpl interface.
func (s *SecurityServiceImpl) GetSecurityUserId(ctx context.Context, req *token.TokenRequest) (resp *token.TokenResponse, err error) {
	// 解析 Token，获取用户 ID
	userId, err := jwt.GetUserId(req.Token)
	if err != nil {
		return &token.TokenResponse{
			Valid: false,
		}, err
	}

	// 返回用户 ID 和 Token 验证状态
	return &token.TokenResponse{
		UserId: *userId,
		Token:  req.Token,
		Valid:  true,
	}, nil
}

// GetSecurityToken implements the SecurityServiceImpl interface.
func (s *SecurityServiceImpl) GetSecurityToken(ctx context.Context, req *token.UserIdRequest) (resp *token.TokenResponse, err error) {
	// 根据用户 ID 生成新的 Token
	newToken, err := jwt.GetToken(req)
	if err != nil {
		return nil, err
	}

	// 返回生成 Token 和相关的用户 ID
	return &token.TokenResponse{
		UserId: req.UserId,
		Token:  newToken,
		Valid:  true,
	}, nil
}
