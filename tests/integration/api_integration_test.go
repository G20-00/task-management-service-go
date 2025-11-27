package integration_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"

	httpdelivery "github.com/G20-00/task-management-service-go/internal/delivery/http"
	repo "github.com/G20-00/task-management-service-go/internal/infrastructure/repository"
	taskusecase "github.com/G20-00/task-management-service-go/internal/usecase/task"
	tasklistusecase "github.com/G20-00/task-management-service-go/internal/usecase/tasklist"
)

// Helper para crear una app Fiber con los handlers reales y repositorios Postgres
func getTestDB(t *testing.T) *sql.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "task_management"
	}
	sqlmode := os.Getenv("DB_SSLMODE")
	if sqlmode == "" {
		sqlmode = "disable"
	}

	connStr := "host=" + host + " port=" + port + " user=" + user +
		" password=" + password + " dbname=" + dbname + " sslmode=" + sqlmode
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	return db
}

func setupApp(t *testing.T) *fiber.App {
	app := fiber.New()
	db := getTestDB(t)
	taskRepo := repo.NewPostgresTaskRepository(db)
	taskListRepo := repo.NewPostgresTaskListRepository(db)
	taskService := taskusecase.NewService(taskRepo)
	taskListService := tasklistusecase.NewService(taskListRepo)
	taskHandler := httpdelivery.NewTaskHandler(taskService)
	taskListHandler := httpdelivery.NewTaskListHandler(taskListService, taskService)

	app.Post("/api/lists", httpdelivery.JWTMiddleware, taskListHandler.CreateTaskList)
	app.Get("/api/lists/:id", httpdelivery.JWTMiddleware, taskListHandler.GetTaskList)
	app.Post("/api/lists/:id/tasks", httpdelivery.JWTMiddleware, taskHandler.CreateTask)
	app.Get("/api/lists/:id/tasks/:taskId", httpdelivery.JWTMiddleware, taskHandler.GetTask)
	app.Patch("/api/lists/:id/tasks/:taskId/state", httpdelivery.JWTMiddleware, taskHandler.UpdateTask)
	app.Delete("/api/lists/:id/tasks/:taskId", httpdelivery.JWTMiddleware, taskHandler.DeleteTask)
	return app
}

func getAuthHeader() string {
	token, err := httpdelivery.GenerateJWT("test-user")
	if err != nil {
		panic("GenerateJWT error: " + err.Error())
	}
	return "Bearer " + token
}

func TestCreateAndGetList_Success(t *testing.T) {
	app := setupApp(t)
	body := `{"name": "Mi lista"}`
	req := httptest.NewRequest(http.MethodPost, "/api/lists", strings.NewReader(body))
	req.Header.Set("Authorization", getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Esperado 201, obtuve %d", resp.StatusCode)
	}
	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if res["id"] == "" {
		t.Error("Se esperaba un id de lista válido")
	}
}

func TestGetList_NotFound(t *testing.T) {
	app := setupApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/lists/no-existe", http.NoBody)
	req.Header.Set("Authorization", getAuthHeader())
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if id, ok := res["id"].(string); ok && id == "no-existe" {
			t.Error("No debería existir una lista con ese id")
		}
	}
}

func TestCreateTaskAndGet_Success(t *testing.T) {
	app := setupApp(t)
	// Crear la lista primero
	listBody := `{"name": "Lista para tarea"}`
	listReq := httptest.NewRequest(http.MethodPost, "/api/lists", strings.NewReader(listBody))
	listReq.Header.Set("Authorization", getAuthHeader())
	listReq.Header.Set("Content-Type", "application/json")
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if listResp.StatusCode != http.StatusCreated {
		t.Fatalf("No se pudo crear la lista, status %d", listResp.StatusCode)
	}
	var listRes map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listRes); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	listID, ok := listRes["id"].(string)
	if !ok {
		t.Fatalf("expected id in listRes")
	}

	body := `{"title": "Tarea 1", "list_id": "` + listID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/lists/"+listID+"/tasks", strings.NewReader(body))
	req.Header.Set("Authorization", getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Esperado 201, obtuve %d", resp.StatusCode)
	}
	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("error decodificando respuesta: %v", err)
	}
	if res["id"] == "" {
		t.Error("Se esperaba un id de tarea válido")
	}
}

func TestGetTask_NotFound(t *testing.T) {
	app := setupApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/lists/list-1/tasks/no-task", http.NoBody)
	req.Header.Set("Authorization", getAuthHeader())
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if res["id"] == "no-task" {
			t.Error("No debería existir una tarea con ese id")
		}
	}
}

func TestChangeTaskState_Success(t *testing.T) {
	app := setupApp(t)
	// Crear la lista primero
	listBody := `{"name": "Lista para cambio estado"}`
	listReq := httptest.NewRequest(http.MethodPost, "/api/lists", strings.NewReader(listBody))
	listReq.Header.Set("Authorization", getAuthHeader())
	listReq.Header.Set("Content-Type", "application/json")
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if listResp.StatusCode != http.StatusCreated {
		t.Fatalf("No se pudo crear la lista, status %d", listResp.StatusCode)
	}
	var listRes map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listRes); err != nil {
		t.Fatalf("error decodificando respuesta: %v", err)
	}
	listIDRaw, ok := listRes["id"].(string)
	if !ok {
		t.Fatalf("el id de la lista no es string: %#v", listRes["id"])
	}
	listID := listIDRaw

	// Crear la tarea primero
	taskBody := `{"title": "Tarea para cambiar estado", "list_id": "` + listID + `"}`
	taskReq := httptest.NewRequest(http.MethodPost, "/api/lists/"+listID+"/tasks", strings.NewReader(taskBody))
	taskReq.Header.Set("Authorization", getAuthHeader())
	taskReq.Header.Set("Content-Type", "application/json")
	taskResp, err := app.Test(taskReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if taskResp.StatusCode != http.StatusCreated {
		t.Fatalf("No se pudo crear la tarea, status %d", taskResp.StatusCode)
	}
	var taskRes map[string]interface{}
	if err := json.NewDecoder(taskResp.Body).Decode(&taskRes); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	taskID, ok := taskRes["id"].(string)
	if !ok {
		t.Fatalf("expected id in taskRes")
	}

	// PATCH con todos los campos requeridos
	patchBody := `{"title": "Tarea para cambiar estado", "status": "completed", "priority": "medium", "list_id": "` + listID + `"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/lists/"+listID+"/tasks/"+taskID+"/state", strings.NewReader(patchBody))
	req.Header.Set("Authorization", getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Esperado 200, obtuve %d", resp.StatusCode)
	}
	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if res["status"] != "completed" {
		t.Errorf("Se esperaba status 'completed', obtuve '%v'", res["status"])
	}
}

func TestChangeTaskState_InvalidStatus(t *testing.T) {
	app := setupApp(t)
	body := `{"status": "invalid"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/lists/list-1/tasks/task-1/state", strings.NewReader(body))
	req.Header.Set("Authorization", getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Error("No debería aceptar un estado inválido")
	}
}

func TestCreateList_EmptyName(t *testing.T) {
	app := setupApp(t)
	body := `{"name": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/lists", strings.NewReader(body))
	req.Header.Set("Authorization", getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode == http.StatusCreated {
		t.Error("No debería crear una lista con nombre vacío")
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	app := setupApp(t)
	// Crear lista válida
	listBody := `{"name": "Lista para tarea vacía"}`
	listReq := httptest.NewRequest(http.MethodPost, "/api/lists", strings.NewReader(listBody))
	listReq.Header.Set("Authorization", getAuthHeader())
	listReq.Header.Set("Content-Type", "application/json")
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	var listRes map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listRes); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	listID, ok := listRes["id"].(string)
	if !ok {
		t.Fatalf("expected id in listRes")
	}
	body := `{"title": "", "list_id": "` + listID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/lists/"+listID+"/tasks", strings.NewReader(body))
	req.Header.Set("Authorization", getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error ejecutando app.Test: %v", err)
	}
	if resp.StatusCode == http.StatusCreated {
		t.Error("No debería crear una tarea con título vacío")
	}
}

func TestCreateTask_InvalidPriority(t *testing.T) {
	app := setupApp(t)
	// Crear lista válida
	listBody := `{"name": "Lista para prioridad"}`
	listReq := httptest.NewRequest(http.MethodPost, "/api/lists", strings.NewReader(listBody))
	listReq.Header.Set("Authorization", getAuthHeader())
	listReq.Header.Set("Content-Type", "application/json")
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	var listRes map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listRes); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	listID, ok := listRes["id"].(string)
	if !ok {
		t.Fatalf("expected id in listRes")
	}
	body := `{"title": "Tarea", "list_id": "` + listID + `", "priority": "super"}`
	req := httptest.NewRequest(http.MethodPost, "/api/lists/"+listID+"/tasks", strings.NewReader(body))
	req.Header.Set("Authorization", getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode == http.StatusCreated {
		t.Error("No debería crear una tarea con prioridad inválida")
	}
}

func TestUpdateTask_InvalidPriority(t *testing.T) {
	app := setupApp(t)
	// Crear lista y tarea válidas
	listBody := `{"name": "Lista para update prioridad"}`
	listReq := httptest.NewRequest(http.MethodPost, "/api/lists", strings.NewReader(listBody))
	listReq.Header.Set("Authorization", getAuthHeader())
	listReq.Header.Set("Content-Type", "application/json")
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	var listRes map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listRes); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	listID, ok := listRes["id"].(string)
	if !ok {
		t.Fatalf("expected id in listRes")
	}
	taskBody := `{"title": "Tarea", "list_id": "` + listID + `"}`
	taskReq := httptest.NewRequest(http.MethodPost, "/api/lists/"+listID+"/tasks", strings.NewReader(taskBody))
	taskReq.Header.Set("Authorization", getAuthHeader())
	taskReq.Header.Set("Content-Type", "application/json")
	taskResp, err := app.Test(taskReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	var taskRes map[string]interface{}
	if err := json.NewDecoder(taskResp.Body).Decode(&taskRes); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	taskID, ok := taskRes["id"].(string)
	if !ok {
		t.Fatalf("expected id in taskRes")
	}
	patchBody := `{"title": "Tarea", "status": "pending", "priority": "super", "list_id": "` + listID + `"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/lists/"+listID+"/tasks/"+taskID+"/state", strings.NewReader(patchBody))
	req.Header.Set("Authorization", getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Error("No debería actualizar una tarea con prioridad inválida")
	}
}

func TestGetTask_NotFound_Explicit(t *testing.T) {
	app := setupApp(t)
	req := httptest.NewRequest(http.MethodGet, "/api/lists/list-1/tasks/no-task", http.NoBody)
	req.Header.Set("Authorization", getAuthHeader())
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Error("No debería encontrar una tarea inexistente")
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	app := setupApp(t)
	patchBody := `{"title": "No existe", "status": "pending", "priority": "medium", "list_id": "list-1"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/lists/list-1/tasks/no-task/state", strings.NewReader(patchBody))
	req.Header.Set("Authorization", getAuthHeader())
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Error("No debería actualizar una tarea inexistente")
	}
}

func TestDeleteTask_AndGet(t *testing.T) {
	app := setupApp(t)
	// Crear lista y tarea válidas
	listBody := `{"name": "Lista para borrar"}`
	listReq := httptest.NewRequest(http.MethodPost, "/api/lists", strings.NewReader(listBody))
	listReq.Header.Set("Authorization", getAuthHeader())
	listReq.Header.Set("Content-Type", "application/json")
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	var listRes map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listRes); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	listID, ok := listRes["id"].(string)
	if !ok {
		t.Fatalf("expected id in listRes")
	}
	taskBody := `{"title": "Tarea a borrar", "list_id": "` + listID + `"}`
	taskReq := httptest.NewRequest(http.MethodPost, "/api/lists/"+listID+"/tasks", strings.NewReader(taskBody))
	taskReq.Header.Set("Authorization", getAuthHeader())
	taskReq.Header.Set("Content-Type", "application/json")
	taskResp, err := app.Test(taskReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	var taskRes map[string]interface{}
	if err := json.NewDecoder(taskResp.Body).Decode(&taskRes); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	taskID, ok := taskRes["id"].(string)
	if !ok {
		t.Fatalf("expected id in taskRes")
	}
	// Eliminar
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/lists/"+listID+"/tasks/"+taskID, http.NoBody)
	deleteReq.Header.Set("Authorization", getAuthHeader())
	resp, err := app.Test(deleteReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Esperado 204 al borrar, obtuve %d", resp.StatusCode)
	}
	// Intentar obtener
	getReq := httptest.NewRequest(http.MethodGet, "/api/lists/"+listID+"/tasks/"+taskID, http.NoBody)
	getReq.Header.Set("Authorization", getAuthHeader())
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if getResp.StatusCode == http.StatusOK {
		t.Error("No debería encontrar la tarea borrada")
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	app := setupApp(t)
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/lists/list-1/tasks/no-task", http.NoBody)
	deleteReq.Header.Set("Authorization", getAuthHeader())
	resp, err := app.Test(deleteReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode == http.StatusNoContent {
		t.Error("No debería borrar una tarea inexistente")
	}
}

func TestCreateTask_WithoutJWT(t *testing.T) {
	app := setupApp(t)
	// Crear lista válida
	listBody := `{"name": "Lista para tarea sin jwt"}`
	listReq := httptest.NewRequest(http.MethodPost, "/api/lists", strings.NewReader(listBody))
	listReq.Header.Set("Authorization", getAuthHeader())
	listReq.Header.Set("Content-Type", "application/json")
	listResp, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	var listRes map[string]interface{}
	if err := json.NewDecoder(listResp.Body).Decode(&listRes); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	listID, ok := listRes["id"].(string)
	if !ok {
		t.Fatalf("expected id in listRes")
	}
	body := `{"title": "Tarea sin jwt", "list_id": "` + listID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/lists/"+listID+"/tasks", strings.NewReader(body))
	// No se agrega Authorization
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode == http.StatusCreated {
		t.Error("No debería crear una tarea sin JWT")
	}
}
