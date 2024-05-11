package usecase

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/LXJ0000/go-backend/domain"
// 	domain_mock "github.com/LXJ0000/go-backend/domain/mock"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"
// )

// func TestRankTopN(t *testing.T) {
// 	tcs := []struct {
// 		name string
// 		mock func(ctrl *gomock.Controller) (domain.InteractionUseCase, domain.PostUsecase)

// 		gotError error
// 		gotPosts []domain.Post
// 	}{
// 		{
// 			name: "Success",
// 			mock: func(ctrl *gomock.Controller) (domain.InteractionUseCase, domain.PostUsecase) {
// 				interaction := domain_mock.NewMockInteractionUseCase(ctrl)
// 				post := domain_mock.NewMockPostUsecase(ctrl)
// 				post.EXPECT().List(gomock.Any(), gomock.Any(), 0, 3).Return([]domain.Post{
// 					{PostID: 1},
// 					{PostID: 2},
// 					{PostID: 3},
// 				}, nil)
// 				interaction.EXPECT().GetByIDs(gomock.Any(), domain.BizPost, []int64{1, 2, 3}).Return(map[int64]domain.Interaction{
// 					1: {BizID: 1, LikeCnt: 3},
// 					2: {BizID: 2, LikeCnt: 2},
// 					3: {BizID: 3, LikeCnt: 1},
// 				}, nil)
// 				post.EXPECT().List(gomock.Any(), gomock.Any(), 3, 3).Return([]domain.Post{}, nil)
// 				interaction.EXPECT().GetByIDs(gomock.Any(), domain.BizPost, nil).Return(map[int64]domain.Interaction{}, nil)
// 				return interaction, post
// 			},
// 			gotPosts: []domain.Post{
// 				{PostID: 1},
// 				{PostID: 2},
// 				{PostID: 3},
// 			},
// 		},
// 	}
// 	for _, tc := range tcs {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			interaction, post := tc.mock(ctrl)
// 			rank := NewPostRankUsecase(interaction, post)
// 			rank.n = 3
// 			rank.batchSize = 3
// 			rank.getScore = func(likeCnt int, updateTime time.Time) float64 {
// 				return float64(likeCnt)
// 			}
// 			posts, err := rank.topN(context.Background())
// 			assert.NoError(t, err)
// 			assert.Equal(t, tc.gotPosts, posts)
// 		})
// 	}
// }
