package logging

import (
	"testing"
	"os"
	"github.com/stretchr/testify/mock"
)

// Fake writer, something that we can set expectations on.
type MockWriter struct{
	mock.Mock
}

// NOTE: This method is not being tested here, code that uses this object is.
func (m *MockWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func Test_UninitializedLogging(*testing.T) {
	// Logger shouldn't crash, even if it is not initialized
	Info("Don't worry")
}

func Test_InitialsedLogging(*testing.T) {
	Init()
	Info("Be happy")
}

func Test_EmptyLogging(*testing.T) {
	Debug("These")
	Info("lines")
	Warning("should")
	Error("go")
	Fatal("nowhere")
}

func Test_UseTargetWhenNotInitialised(*testing.T) {
	AddTarget(os.Stdout, DebugLevel)
	Debug("True Survivor")
}

func Test_UseTargetWhenInitialised(*testing.T) {
	Init()
	AddTarget(os.Stdout, DebugLevel)
	Info("Smooth operator")
}

func Test_MultipleTargetsWhenNotInitialised(*testing.T) {
	AddTarget(os.Stdout, DebugLevel)
	AddTarget(os.Stderr, ErrorLevel)
	Info("Bear Grylls")
}

func Test_CallOnlyRelevantTargets(t *testing.T) {
	mockTargetA := new(MockWriter)
	mockTargetB := new(MockWriter)
	mockTargetC := new(MockWriter)

	// Return value doesn't matter as we don't use it now anywhere
	mockTargetA.On("Write", mock.Anything).Return(33, nil)
	mockTargetC.On("Write", mock.Anything).Return(33, nil)

	Init()
	AddTarget(mockTargetA, DebugLevel)
	AddTarget(mockTargetB, ErrorLevel)
	AddTarget(mockTargetC, InfoLevel)

	Info("Useless, but true information.")

	mockTargetA.AssertExpectations(t)
	mockTargetB.AssertNotCalled(t, "Write", mock.Anything)
	mockTargetC.AssertCalled(t, "Write", mock.Anything)
}

func Test_NullTargetsAreNotWelcome(t *testing.T) {
	// Shouldn't crash on passing nil as target
	Init()
	AddTarget(nil, DebugLevel)
	Debug("You can't stop me!")
}
