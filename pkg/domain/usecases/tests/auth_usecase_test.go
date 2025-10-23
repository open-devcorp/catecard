package usecases

/*
func TestSignUp(t *testing.T) {
	// GIVEN
	logger := log.New(io.Discard, "", 0)
	mockRepo := repositories.NewMockUserRepository()
	uc := usecases.NewAuthUseCase(logger, mockRepo)

	input := usecases.SignupStruct{
		Username: "newuser@example.com",
		Password: "pass123",
		Role:     entities.CATECHIST,
	}

	// WHEN
	user, err := uc.SignUp(input)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, input.Username, user.Username)
	assert.Equal(t, input.Password, user.Password)
	assert.Equal(t, input.Role, user.Role)
}

func TestLogin(t *testing.T) {
	// GIVEN
	logger := log.New(io.Discard, "", 0)
	mockRepo := repositories.NewMockUserRepository()
	uc := usecases.NewAuthUseCase(logger, mockRepo)

	signupInput := usecases.SignupStruct{
		Username: "existinguser@example.com",
		Password: "secret123",
		Role:     entities.CATECHIST,
	}

	// ensure user exists
	_, err := uc.SignUp(signupInput)
	assert.Nil(t, err)

	loginInput := usecases.LoginStruct{
		Username: signupInput.Username,
		Password: signupInput.Password,
	}

	// WHEN
	user, err := uc.Login(loginInput)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, signupInput.Username, user.Username)
	assert.Equal(t, signupInput.Password, user.Password)
	assert.Equal(t, signupInput.Role, user.Role)
}

*/
