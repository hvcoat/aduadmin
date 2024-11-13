package service

import (
	"context"

	v1 "hd/api/helloworld/v1"
	"hd/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	_, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	htmlContent := `
<!DOCTYPE html>
<html>

<body>
    <form action="your_action_page.php" method="post">
        <input type="image" src="your_image_path.jpg" alt="Submit" width="100" height="50">
    </form>
</body>

</html>
`
	return &v1.HelloReply{Message: htmlContent}, nil
}
