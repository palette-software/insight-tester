package logging

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Fake writer, something that we can set expectations on.
type MockWriter struct {
	mock.Mock
}

// NOTE: This method is not being tested here, code that uses this object is.
func (m *MockWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

// Test suite so that logger can start from zero in all tests
type LoggerTestSuite struct {
	suite.Suite
}

// Revert logger to "zero" state, before all tests
func (suite *LoggerTestSuite) SetupTest() {
	loggers = nil
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func Test_LoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (suite *LoggerTestSuite) Test_EmptyLogging() {
	Debug("These")
	Info("lines")
	Warning("should")
	Error("go")
	Fatal("nowhere")
}

func (suite *LoggerTestSuite) Test_SingleTarget() {
	AddTarget(os.Stdout, LevelDebug)
	Info("Smooth operator")
}

func (suite *LoggerTestSuite) Test_MultipleTargetsButTriggerOnlyOne() {
	AddTarget(os.Stdout, LevelDebug)
	AddTarget(os.Stderr, LevelError)
	Info("One ring above all")
}

func (suite *LoggerTestSuite) Test_MultipleTargetsButTriggerBoth() {
	AddTarget(os.Stdout, LevelDebug)
	AddTarget(os.Stderr, LevelError)
	Error("Get together and feel alright")
}

func (suite *LoggerTestSuite) Test_CallOnlyRelevantTargets() {
	mockTargetA := new(MockWriter)
	mockTargetB := new(MockWriter)
	mockTargetC := new(MockWriter)

	// Return value doesn't matter as we don't use it now anywhere
	mockTargetA.On("Write", mock.Anything).Return(33, nil)
	mockTargetC.On("Write", mock.Anything).Return(33, nil)

	AddTarget(mockTargetA, LevelDebug)
	AddTarget(mockTargetB, LevelError)
	AddTarget(mockTargetC, LevelInfo)

	Info("Useless, but true information.")

	mockTargetA.AssertExpectations(suite.T())
	mockTargetB.AssertNotCalled(suite.T(), "Write", mock.Anything)
	mockTargetC.AssertCalled(suite.T(), "Write", mock.Anything)
}

func (suite *LoggerTestSuite) Test_NullTargetsAreNotWelcome() {
	// Shouldn't crash on passing nil as target
	AddTarget(nil, LevelDebug)
	Debug("You can't stop me!")
}
