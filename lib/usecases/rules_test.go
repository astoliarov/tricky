package usecases

import (
	"tricky/lib/domain"

	"errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"fmt"
	"net/url"
	"testing"
)

type SameRuleMatcher struct {
	rule *domain.RedirectionRule
}

func (m *SameRuleMatcher) Matches(x interface{}) bool {
	rule, ok := x.(domain.RedirectionRule)
	if !ok {
		fmt.Println("not casted")
		return false
	}

	if rule.Key == m.rule.Key && rule.Url == m.rule.Url {
		fmt.Println("not satisfied")
		return true
	}

	return false
}

func (m *SameRuleMatcher) String() string {
	return "SameRuleMatcher"
}

type RulesInteractorTestSuite struct {
	suite.Suite
	interactor *RulesUseCase
	repository *MockIRulesDAO
}

func (suite *RulesInteractorTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.repository = NewMockIRulesDAO(mockCtrl)
	suite.interactor = &RulesUseCase{rulesDAO: suite.repository}
}

func (suite *RulesInteractorTestSuite) getTestRule() *domain.RedirectionRule {
	u, _ := url.Parse("http://testdomain.com")
	return &domain.RedirectionRule{
		Key: "test",
		Url: u,
	}
}

func (suite *RulesInteractorTestSuite) TestFindUrlSuccess() {
	rule := suite.getTestRule()
	keys := []string{rule.Key}

	suite.repository.EXPECT().GetByKey(rule.Key).Return(rule, nil)

	foundedUrl, err := suite.interactor.FindUrl(keys)

	assert.Equal(suite.T(), foundedUrl, rule.Url)
	assert.Nil(suite.T(), err)
}

func (suite *RulesInteractorTestSuite) TestFindUrlStorageError() {
	key := "test"
	e := errors.New("Test error")

	keys := []string{key}

	suite.repository.EXPECT().GetByKey(key).Return(nil, e)

	foundedUrl, err := suite.interactor.FindUrl(keys)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), err, e)
	assert.Nil(suite.T(), foundedUrl)
}

func (suite *RulesInteractorTestSuite) TestFindUrlSecondKeySuccess() {
	rule := suite.getTestRule()
	keys := []string{"empty", rule.Key}

	suite.repository.EXPECT().GetByKey("empty").Return(nil, nil)
	suite.repository.EXPECT().GetByKey(rule.Key).Return(rule, nil)

	foundedUrl, _ := suite.interactor.FindUrl(keys)

	assert.Equal(suite.T(), foundedUrl, rule.Url)
}

func (suite *RulesInteractorTestSuite) TestFindUrlNoRule() {
	key := "test"
	keys := []string{key}

	suite.repository.EXPECT().GetByKey(key).Return(nil, nil)

	foundedUrl, _ := suite.interactor.FindUrl(keys)

	assert.Nil(suite.T(), foundedUrl)
}

//func (suite *RulesInteractorTestSuite) TestAddRuleSuccess() {
//	rule := suite.getTestRule()
//
//	suite.repository.EXPECT().Add(SameRuleMatcher{rule})
//
//	err := suite.interactor.AddRule(rule.Key, rule.Url)
//
//	assert.Nil(suite.T(), err)
//}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(RulesInteractorTestSuite))
}
