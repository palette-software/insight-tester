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

// Fake os.Exit(int) function
type MockOs struct {
	mock.Mock
}

// NOTE: This method is not being tested here, code that uses this object is.
func (m *MockOs) Exit(code int) {
	m.Called(code)
	return
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

func (suite *LoggerTestSuite) Test_NullTargetsAreNotWelcome() {
	// Shouldn't crash on passing nil as target
	AddTarget(nil, LevelDebug)
	Debug("You can't stop me!")
}

// Test suite so that logger can start from zero in all tests
type LoggerMockSuite struct {
	suite.Suite
	mockTargetA *MockWriter
	mockTargetB *MockWriter
	mockTargetC *MockWriter
	mockOs *MockOs
}

// Initialize some mocks
func (suite *LoggerMockSuite) SetupTest() {
	loggers = nil

	suite.mockTargetA = new(MockWriter)
	suite.mockTargetB = new(MockWriter)
	suite.mockTargetC = new(MockWriter)

	suite.mockOs = new(MockOs)

	// Patch, so that testing process won't quit on Fatal log tests
	osExit = suite.mockOs.Exit
}

func (suite *LoggerMockSuite) TearDownTest() {
	// Revert back to the original os.Exit() function
	osExit = os.Exit

	suite.mockTargetA.AssertExpectations(suite.T())
	suite.mockTargetB.AssertExpectations(suite.T())
	suite.mockTargetC.AssertExpectations(suite.T())
	suite.mockOs.AssertExpectations(suite.T())
}

func (suite *LoggerMockSuite) Test_CallOnlyRelevantTargets() {
	// Return value doesn't matter as we don't use it now anywhere
	suite.mockTargetA.On("Write", mock.Anything).Return(33, nil)
	suite.mockTargetC.On("Write", mock.Anything).Return(33, nil)

	AddTarget(suite.mockTargetA, LevelDebug)
	AddTarget(suite.mockTargetB, LevelError)
	AddTarget(suite.mockTargetC, LevelInfo)

	Info("Useless, but true information.")
}

func (suite *LoggerMockSuite) Test_Fatal() {
	// Return value doesn't matter as we don't use it now anywhere
	suite.mockTargetA.On("Write", mock.Anything).Return(33, nil)
	suite.mockTargetB.On("Write", mock.Anything).Return(33, nil)
	suite.mockTargetC.On("Write", mock.Anything).Return(33, nil)
	suite.mockOs.On("Exit", 1)

	AddTarget(suite.mockTargetA, LevelDebug)
	AddTarget(suite.mockTargetB, LevelError)
	AddTarget(suite.mockTargetC, LevelInfo)

	Fatal("Run, everyone, run")

	suite.mockOs.AssertNumberOfCalls(suite.T(), "Exit", 1)
}

func (suite *LoggerMockSuite) Test_Fatalf() {
	// Return value doesn't matter as we don't use it now anywhere
	suite.mockTargetA.On("Write", mock.Anything).Return(33, nil)
	suite.mockTargetB.On("Write", mock.Anything).Return(33, nil)
	suite.mockTargetC.On("Write", mock.Anything).Return(33, nil)
	suite.mockOs.On("Exit", 1)

	AddTarget(suite.mockTargetA, LevelDebug)
	AddTarget(suite.mockTargetB, LevelError)
	AddTarget(suite.mockTargetC, LevelInfo)

	Fatalf("The number of the week: %d", 7)

	suite.mockOs.AssertNumberOfCalls(suite.T(), "Exit", 1)
}

// Now run the mock suite
func Test_LoggerMockSuite(t *testing.T) {
	suite.Run(t, new(LoggerMockSuite))
}
