TEST_DIR=.
COVERAGE_TMP=coverage.out.tmp
COVERAGE_OUT=coverage.out
FILES_TO_CLEAN=*.out *.out.tmp *DS_Store 
MOCKS="mocks"

test:
	# @echo "Делаем моки..."
	# mockgen -source= -destination= -package=
	@echo "Запуск тестов..."
	go test -v -race -coverpkg=./... -coverprofile=$(COVERAGE_TMP) $(TEST_DIR)/...
	@echo "Обработка покрытия..."

	# Добавляем условие для исключения файлов
	# cat $(COVERAGE_TMP) | grep -vE '' > $(COVERAGE_OUT)
	cat $(COVERAGE_TMP) > $(COVERAGE_OUT)
	rm $(COVERAGE_TMP)

	go tool cover -func=$(COVERAGE_OUT)

	@echo "Удаление временных файлов..."
	rm -f $(FILES_TO_CLEAN)
	rm -rf $(MOCKS)
	@echo "Тесты завершены"
