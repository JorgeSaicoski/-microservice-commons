// database/migration.go
package database

import (
	"fmt"

	"github.com/JorgeSaicoski/pgconnect"
)

// MigrationOptions holds migration configuration
type MigrationOptions struct {
	DropTables    bool // Drop tables before migration (dangerous!)
	CreateIndexes bool // Create indexes after migration
	Verbose       bool // Print migration details
}

// DefaultMigrationOptions returns safe default migration options
func DefaultMigrationOptions() MigrationOptions {
	return MigrationOptions{
		DropTables:    false,
		CreateIndexes: true,
		Verbose:       true,
	}
}

// Migrator handles database migrations
type Migrator struct {
	db      *pgconnect.DB
	options MigrationOptions
	models  []interface{}
}

// NewMigrator creates a new database migrator
func NewMigrator(db *pgconnect.DB, options MigrationOptions) *Migrator {
	return &Migrator{
		db:      db,
		options: options,
		models:  make([]interface{}, 0),
	}
}

// AddModels adds models to be migrated
func (m *Migrator) AddModels(models ...interface{}) *Migrator {
	m.models = append(m.models, models...)
	return m
}

// Migrate runs the migration
func (m *Migrator) Migrate() error {
	if len(m.models) == 0 {
		return fmt.Errorf("no models to migrate")
	}

	if m.options.Verbose {
		fmt.Printf("Starting migration for %d models...\n", len(m.models))
	}

	// Drop tables if requested (use with extreme caution!)
	if m.options.DropTables {
		if m.options.Verbose {
			fmt.Println("WARNING: Dropping existing tables...")
		}
		for _, model := range m.models {
			if err := m.db.DB.Migrator().DropTable(model); err != nil {
				fmt.Printf("Warning: Failed to drop table for model %T: %v\n", model, err)
			}
		}
	}

	// Run auto migration
	if err := m.db.AutoMigrate(m.models...); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	if m.options.Verbose {
		fmt.Println("Migration completed successfully")
	}

	// Create indexes if requested
	if m.options.CreateIndexes {
		if err := m.createIndexes(); err != nil {
			fmt.Printf("Warning: Failed to create some indexes: %v\n", err)
		}
	}

	return nil
}

// createIndexes creates common indexes for better performance
func (m *Migrator) createIndexes() error {
	// Add common indexes here based on your models
	// This is where you'd add indexes that are common across services

	if m.options.Verbose {
		fmt.Println("Creating additional indexes...")
	}

	// Example indexes - customize based on your needs
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_created_at ON base_projects(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_updated_at ON base_projects(updated_at);",
		"CREATE INDEX IF NOT EXISTS idx_owner_id ON base_projects(owner_id);",
		"CREATE INDEX IF NOT EXISTS idx_company_id ON base_projects(company_id);",
		"CREATE INDEX IF NOT EXISTS idx_status ON base_projects(status);",
		"CREATE INDEX IF NOT EXISTS idx_company_members_user ON company_members(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_company_members_company ON company_members(company_id);",
		"CREATE INDEX IF NOT EXISTS idx_company_members_status ON company_members(status);",
		"CREATE INDEX IF NOT EXISTS idx_project_members_user ON project_members(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_project_members_project ON project_members(project_id);",
	}

	for _, indexSQL := range indexes {
		if err := m.db.DB.Exec(indexSQL).Error; err != nil {
			fmt.Printf("Warning: Failed to create index: %s - %v\n", indexSQL, err)
		}
	}

	return nil
}

// QuickMigrate is a convenience function for simple migrations
func QuickMigrate(db *pgconnect.DB, models ...interface{}) error {
	migrator := NewMigrator(db, DefaultMigrationOptions())
	return migrator.AddModels(models...).Migrate()
}

// UnsafeMigrate drops all tables and recreates them (DANGEROUS!)
func UnsafeMigrate(db *pgconnect.DB, models ...interface{}) error {
	options := MigrationOptions{
		DropTables:    true,
		CreateIndexes: true,
		Verbose:       true,
	}

	fmt.Println("⚠️  WARNING: This will DROP ALL TABLES! ⚠️")

	migrator := NewMigrator(db, options)
	return migrator.AddModels(models...).Migrate()
}
