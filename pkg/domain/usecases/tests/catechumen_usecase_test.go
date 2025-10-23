package usecases

/*
func TestAdd(t *testing.T) {
	// GIVEN
	logger := log.New(nil, "", 0)
	catechumenRepo := repositories.NewMockCatechumenRepository()
	qrRepo := repositories.NewMockQrRepository()
	fakeCatechumen := entities.FakeCatechumen()
	uc := usecases.NewCatechumenUsecase(logger, catechumenRepo, qrRepo)

	// WHEN
	catechumen, qr, err := uc.Add(&fakeCatechumen)

	// THEN
	assert.Nil(t, err)

	assert.NotNil(t, catechumen)
	assert.Equal(t, 1, catechumen.ID)

	assert.NotNil(t, qr)
	assert.Equal(t, 1, len(catechumenRepo.Catechumens))
	assert.Equal(t, 1, len(qrRepo.Qrs))
}

func TestUpdate(t *testing.T) {
	// GIVEN
	logger := log.New(nil, "", 0)
	catechumenRepo := repositories.NewMockCatechumenRepository()
	qrRepo := repositories.NewMockQrRepository()
	fakeCatechumen := entities.FakeCatechumen()

	catechumenRepo.Add(&fakeCatechumen)
	fakeCatechumen.FullName = "Updated Name"

	uc := usecases.NewCatechumenUsecase(logger, catechumenRepo, qrRepo)

	// WHEN
	updatedCatechumen, err := uc.Update(&fakeCatechumen)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, updatedCatechumen)
	assert.Equal(t, 1, updatedCatechumen.ID)

	assert.Equal(t, "Updated Name", updatedCatechumen.FullName)
}

func TestGetAll(t *testing.T) {
	// GIVEN
	logger := log.New(nil, "", 0)
	catechumenRepo := repositories.NewMockCatechumenRepository()
	qrRepo := repositories.NewMockQrRepository()
	fakeCatechumen := entities.FakeCatechumen()
	uc := usecases.NewCatechumenUsecase(logger, catechumenRepo, qrRepo)

	// WHEN
	uc.Add(&fakeCatechumen)
	catechumens, err := uc.GetAll()

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, 1, len(catechumens))
}


func TestGetById(t *testing.T) {
	// GIVEN
	logger := log.New(nil, "", 0)
	catechumenRepo := repositories.NewMockCatechumenRepository()
	qrRepo := repositories.NewMockQrRepository()
	fakeCatechumen := entities.FakeCatechumen()
	uc := usecases.NewCatechumenUsecase(logger, catechumenRepo, qrRepo)

	// WHEN
	uc.Add(&fakeCatechumen)
	catechumen, err := uc.GetById(1)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, catechumen)
	assert.Equal(t, 1, catechumen.ID)
}
*/
