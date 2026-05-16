package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	appMiddleware "pos-desktop-tauri/backend/internal/middleware"
	"pos-desktop-tauri/backend/internal/handler"
)

func New() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	// Internal (discovery — tanpa auth, hanya dari localhost)
	r.Post("/api/internal/discovery/advertise", handler.DiscoveryAdvertise)
	r.Get("/api/internal/discovery/scan",       handler.DiscoveryScan)
	r.Post("/api/internal/discovery/stop",      handler.DiscoveryStop)

	// Public
	r.Post("/api/auth/login", handler.Login)

	// Auth required
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.Auth)

		r.Get("/api/auth/me", handler.Me)
		r.Get("/api/permissions", handler.PermissionIndex)

		r.Get("/api/toko", handler.TokoShow)
		r.Put("/api/toko", handler.TokoUpdate)

		// Produk
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Can("produk"))
			r.Get("/api/produk", handler.ProdukIndex)
			r.Post("/api/produk", handler.ProdukStore)
			r.Get("/api/produk/trashed", handler.ProdukTrashed)
			r.Get("/api/produk/next-kode", handler.ProdukNextKode)
			r.Get("/api/produk/{id}", handler.ProdukShow)
			r.Put("/api/produk/{id}", handler.ProdukUpdate)
			r.Delete("/api/produk/{id}", handler.ProdukDestroy)
			r.Post("/api/produk/{id}/restore", handler.ProdukRestore)
			r.Get("/api/produk/{id}/satuan", handler.ProdukGetSatuan)
			r.Post("/api/produk/{id}/satuan", handler.ProdukStoreSatuan)
			r.Put("/api/produk/{id}/satuan/{satuanId}", handler.ProdukUpdateSatuan)
			r.Delete("/api/produk/{id}/satuan/{satuanId}", handler.ProdukDestroySatuan)

			r.Get("/api/kategori", handler.KategoriIndex)
			r.Post("/api/kategori", handler.KategoriStore)
			r.Put("/api/kategori/{id}", handler.KategoriUpdate)
			r.Delete("/api/kategori/{id}", handler.KategoriDestroy)

			r.Get("/api/pemasok", handler.PemasokIndex)
			r.Post("/api/pemasok", handler.PemasokStore)
			r.Put("/api/pemasok/{id}", handler.PemasokUpdate)
			r.Delete("/api/pemasok/{id}", handler.PemasokDestroy)
		})

		// Pengaturan
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Can("pengaturan"))
			r.Get("/api/users", handler.UserIndex)
			r.Post("/api/users", handler.UserStore)
			r.Put("/api/users/{id}", handler.UserUpdate)
			r.Delete("/api/users/{id}", handler.UserDestroy)

			r.Get("/api/roles", handler.RoleIndex)
			r.Post("/api/roles", handler.RoleStore)
			r.Put("/api/roles/{id}", handler.RoleUpdate)
		})
	})

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
