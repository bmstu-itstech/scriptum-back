package service_test

import "testing"

func TestRepositories(t *testing.T) {
	t.Run("FileRepository", testFileRepository)
	t.Run("ScriptRepository", testScriptRepository)

	t.Run("JobRepository", testJobRepository)

}
