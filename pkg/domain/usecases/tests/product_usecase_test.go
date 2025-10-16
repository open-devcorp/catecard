package usecases

// func TestAddProduct(t *testing.T) {

// 	//GIVEN
// 	logger := log.New(nil, "", 0)
// 	mockRepo := repositories.NewMockProductRepository()
// 	fakeProduct := entities.FakeProduct()
// 	fakeProduct.Name = "PRUEBA TEST"
// 	uc := usecases.NewProductUsecase(logger, mockRepo)
// 	//WHEN

// 	err := uc.Add(&fakeProduct)

// 	//THEN

// 	assert.Nil(t, err)
// 	assert.Equal(t, len(mockRepo.Products), 1)

// 	log.Printf("LIST: %v", mockRepo.Products[0])

// }

// func TestGetAll(t *testing.T) {
// 	//GIVEN
// 	logger := log.New(nil, "", 0)
// 	mockRepo := repositories.NewMockProductRepository()
// 	uc := usecases.NewProductUsecase(logger, mockRepo)
// 	fake := entities.FakeProduct()
// 	//WHEN
// 	uc.Add(&fake) //ADD 1
// 	products, err := uc.GetAll()

// 	assert.Nil(t, err)
// 	assert.Equal(t, len(products), 1)

// }
