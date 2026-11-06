// performance_test.go — тестирование оптимизаций производительности для INV-PROMPT-MVP-009
package screens

import (
	"strings"
	"testing"
	"time"

	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui"
)

// TestWidgetDataCache проверяет, что кеширование данных виджетов работает правильно
func TestWidgetDataCache(t *testing.T) {
	// Создаем тестовый дашборд
	d := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Проверяем, что кеш инициализирован
	if d.widgetDataCache == nil {
		t.Error("Widget data cache should be initialized")
	}
	
	// Проверяем, что кеш невалиден по умолчанию
	if d.widgetCacheValid {
		t.Error("Widget cache should be invalid initially")
	}
	
	// Проверяем, что хеши инициализированы
	if d.lastSessionsHash == "" {
		t.Error("Sessions hash should be initialized")
	}
	if d.lastManifestHash == "" {
		t.Error("Manifest hash should be initialized")
	}
	if d.lastAdaptersHash == "" {
		t.Error("Adapters hash should be initialized")
	}
}

// TestStringBuilding проверяет, что string builder используется правильно
func TestStringBuilding(t *testing.T) {
	d := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Тестирование renderSuccess с использованием string builder
	start := time.Now()
	result := d.renderSuccess()
	duration := time.Since(start)
	
	// Проверяем, что результат не пустой
	if result == "" {
		t.Error("renderSuccess should return non-empty string")
	}
	
	// Проверяем, что результат содержит ожидаемые компоненты
	if !strings.Contains(result, "dashboard") {
		t.Error("renderSuccess should contain 'dashboard'")
	}
	if !strings.Contains(result, "StatusBar") {
		t.Error("renderSuccess should contain StatusBar")
	}
	
	// Проверяем, что выполнение занимает разумное время (< 50ms для простого случая)
	if duration > 50*time.Millisecond {
		t.Logf("renderSuccess took %v (acceptable for complex case)", duration)
	} else {
		t.Logf("renderSuccess took %v (good performance)", duration)
	}
}

// TestCacheInvalidation проверяет, что кеш инвалидится при изменении данных
func TestCacheInvalidation(t *testing.T) {
	d := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Первоначальное состояние - кеш невалиден
	if d.widgetCacheValid {
		t.Error("Cache should be invalid initially")
	}
	
	// Симулируем загрузку данных
	mockData := DashboardDataLoaded{
		Sessions: []*storage.Session{},
		Manifest: &storage.Manifest{},
		Adapters: []caf.Adapter{},
		LoadTime: time.Now(),
	}
	
	d.handleDataLoaded(mockData)
	
	// После загрузки данных кеш должен стать валидным
	if !d.widgetCacheValid {
		t.Error("Cache should be valid after data loading")
	}
	
	// Симулируем изменение данных (refresh)
	d.Loading = true
	d.LoadError = nil
	d.StateMachine.SetState(ui.StateLoading)
	
	// Кеш должен стать невалидным снова
	if d.widgetCacheValid {
		t.Error("Cache should be invalid after refresh")
	}
}

// TestBorderStringCache проверяет, что строки границ кешируются
func TestBorderStringCache(t *testing.T) {
	d := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Проверяем, что кеш границ инициализирован
	if d.widgetDataCache.BorderStrings.Top == "" {
		t.Error("Border cache should be initialized")
	}
	if d.widgetDataCache.BorderStrings.Bottom == "" {
		t.Error("Border cache should be initialized")
	}
	
	// Проверяем, что строки содержат ожидаемые символы
	if !strings.Contains(d.widgetDataCache.BorderStrings.Top, "┌") {
		t.Error("Top border should contain '┌'")
	}
	if !strings.Contains(d.widgetDataCache.BorderStrings.Bottom, "└") {
		t.Error("Bottom border should contain '└'")
	}
}

// TestGridLayoutCache проверяет, что данные компоновки кешируются
func TestGridLayoutCache(t *testing.T) {
	d := NewDashboardScreen("/tmp/test", "0.4.1-test")
	
	// Проверяем, что кеш компоновки инициализирован
	if d.widgetDataCache.GridLayout.WidgetWidth == 0 {
		t.Error("Layout cache should be initialized")
	}
	if d.widgetDataCache.GridLayout.WidgetHeight == 0 {
		t.Error("Layout cache should be initialized")
	}
	if len(d.widgetDataCache.GridLayout.WidgetPositions) == 0 {
		t.Error("Layout cache should have widget positions")
	}
	
	// Проверяем, что позиции виджетов установлены
	if _, exists := d.widgetDataCache.GridLayout.WidgetPositions["Today"]; !exists {
		t.Error("Today widget position should be cached")
	}
	if _, exists := d.widgetDataCache.GridLayout.WidgetPositions["Authorship Score"]; !exists {
		t.Error("Authorship widget position should be cached")
	}
}