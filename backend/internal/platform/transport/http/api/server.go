package api

import (
	"context"
	"messenger/messenger/internal/platform/commandbus"
	"messenger/messenger/internal/platform/querybus"
)

var _ StrictServerInterface = (*Server)(nil) 

type Server struct {
	commandBus *commandbus.Bus
	queryBus   *querybus.Bus
}

func NewServer(
	commandBus *commandbus.Bus,
	queryBus   *querybus.Bus,
) *Server {
	return &Server{
		commandBus: commandBus,
		queryBus: queryBus,
	}
}

func (s *Server) Register(ctx context.Context, request RegisterRequestObject) (RegisterResponseObject, error) {
	
	return nil, nil
}

func (s *Server) Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error) {
	return nil, nil
}

func (s *Server) RefreshTokens(ctx context.Context, request RefreshTokensRequestObject) (RefreshTokensResponseObject, error) {
	return nil, nil
}

func (s *Server) GetMe(ctx context.Context, request GetMeRequestObject) (GetMeResponseObject, error) {
	return nil, nil
}

func (s *Server) GetMyChatList(ctx context.Context, request GetMyChatListRequestObject) (GetMyChatListResponseObject, error) {
	return nil, nil
}

