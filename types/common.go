package types

import (
	"time"
)

// BaseModel represents common fields for database models
type BaseModel struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// UUIDBaseModel represents common fields with UUID primary key
type UUIDBaseModel struct {
	ID        string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// AuditableModel includes audit fields
type AuditableModel struct {
	BaseModel
	CreatedBy *uint `json:"created_by,omitempty" gorm:"index"`
	UpdatedBy *uint `json:"updated_by,omitempty" gorm:"index"`
}

// Status represents common status types
type Status string

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusPending   Status = "pending"
	StatusApproved  Status = "approved"
	StatusRejected  Status = "rejected"
	StatusCancelled Status = "cancelled"
	StatusCompleted Status = "completed"
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

// Priority represents priority levels
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// Role represents user roles
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
	RoleGuest     Role = "guest"
	RoleOwner     Role = "owner"
	RoleMember    Role = "member"
	RoleViewer    Role = "viewer"
	RoleEditor    Role = "editor"
)

// Permission represents permissions
type Permission string

const (
	PermissionRead   Permission = "read"
	PermissionWrite  Permission = "write"
	PermissionDelete Permission = "delete"
	PermissionAdmin  Permission = "admin"
	PermissionCreate Permission = "create"
	PermissionUpdate Permission = "update"
)

// Environment represents deployment environments
type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
	EnvTest        Environment = "test"
)

// LogLevel represents logging levels
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// Gender represents gender options
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// Currency represents currency codes
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencyJPY Currency = "JPY"
	CurrencyCAD Currency = "CAD"
	CurrencyAUD Currency = "AUD"
)

// Language represents language codes
type Language string

const (
	LanguageEnglish    Language = "en"
	LanguageSpanish    Language = "es"
	LanguageFrench     Language = "fr"
	LanguageGerman     Language = "de"
	LanguageItalian    Language = "it"
	LanguagePortuguese Language = "pt"
	LanguageJapanese   Language = "ja"
	LanguageChinese    Language = "zh"
)

// UserStatus represents user account statuses
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
	UserStatusPending   UserStatus = "pending"
	UserStatusVerified  UserStatus = "verified"
)

// NotificationType represents notification types
type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeError   NotificationType = "error"
	NotificationTypeSuccess NotificationType = "success"
)

// MediaType represents file/media types
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeDocument MediaType = "document"
	MediaTypeArchive  MediaType = "archive"
	MediaTypeOther    MediaType = "other"
)

// SortOrder represents sorting directions
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// SortBy represents common sortable fields
type SortBy string

const (
	SortByID        SortBy = "id"
	SortByName      SortBy = "name"
	SortByCreatedAt SortBy = "created_at"
	SortByUpdatedAt SortBy = "updated_at"
	SortByTitle     SortBy = "title"
	SortByPriority  SortBy = "priority"
	SortByStatus    SortBy = "status"
)

// ContactInfo represents contact information
type ContactInfo struct {
	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
	Website string `json:"website,omitempty"`
}

// Address represents a physical address
type Address struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
}

// Coordinates represents geographical coordinates
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Money represents monetary amounts
type Money struct {
	Amount   int64    `json:"amount"` // Amount in smallest currency unit (cents)
	Currency Currency `json:"currency"`
}

// DateRange represents a date range
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// TimeRange represents a time range within a day
type TimeRange struct {
	StartTime string `json:"start_time"` // HH:MM format
	EndTime   string `json:"end_time"`   // HH:MM format
}

// Metadata represents flexible key-value metadata
type Metadata map[string]interface{}

// Settings represents application settings
type Settings map[string]interface{}

// Tags represents a list of tags
type Tags []string

// KeyValue represents a key-value pair
type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Option represents a selectable option
type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// File represents file information
type File struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	URL      string `json:"url,omitempty"`
	Path     string `json:"path,omitempty"`
	Checksum string `json:"checksum,omitempty"`
}

// Image represents image information
type Image struct {
	File
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Alt    string `json:"alt,omitempty"`
}

// Link represents a hyperlink
type Link struct {
	URL   string `json:"url"`
	Text  string `json:"text,omitempty"`
	Title string `json:"title,omitempty"`
}

// Social represents social media links
type Social struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
	Username string `json:"username,omitempty"`
}

// Notification represents a notification
type Notification struct {
	ID        string           `json:"id"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Read      bool             `json:"read"`
	CreatedAt time.Time        `json:"created_at"`
	Data      Metadata         `json:"data,omitempty"`
}

// APIKey represents an API key
type APIKey struct {
	BaseModel
	Name        string     `json:"name" gorm:"not null"`
	Key         string     `json:"key" gorm:"uniqueIndex;not null"`
	UserID      uint       `json:"user_id" gorm:"not null;index"`
	Permissions []string   `json:"permissions" gorm:"type:text[]"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	BaseModel
	UserID     *uint    `json:"user_id,omitempty" gorm:"index"`
	Action     string   `json:"action" gorm:"not null"`
	Resource   string   `json:"resource" gorm:"not null"`
	ResourceID string   `json:"resource_id,omitempty" gorm:"index"`
	OldData    Metadata `json:"old_data,omitempty" gorm:"type:jsonb"`
	NewData    Metadata `json:"new_data,omitempty" gorm:"type:jsonb"`
	IPAddress  string   `json:"ip_address,omitempty"`
	UserAgent  string   `json:"user_agent,omitempty"`
}

// IsValid checks if a status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusPending, StatusApproved,
		StatusRejected, StatusCancelled, StatusCompleted, StatusDraft,
		StatusPublished, StatusArchived:
		return true
	default:
		return false
	}
}

// IsValid checks if a priority is valid
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	default:
		return false
	}
}

// IsValid checks if a role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleModerator, RoleGuest,
		RoleOwner, RoleMember, RoleViewer, RoleEditor:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (s Status) String() string {
	return string(s)
}

// String returns the string representation
func (p Priority) String() string {
	return string(p)
}

// String returns the string representation
func (r Role) String() string {
	return string(r)
}
