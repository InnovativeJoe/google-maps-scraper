package web

import "strings"

// FieldSelection controls which optional Google Maps fields are emitted in the
// web UI CSV export.
type FieldSelection struct {
	Link             bool `json:"link"`
	Title            bool `json:"title"`
	Category         bool `json:"category"`
	Address          bool `json:"address"`
	OpenHours        bool `json:"open_hours"`
	Website          bool `json:"website"`
	Phone            bool `json:"phone"`
	PlusCode         bool `json:"plus_code"`
	ReviewCount      bool `json:"review_count"`
	ReviewRating     bool `json:"review_rating"`
	ReviewsPerRating bool `json:"reviews_per_rating"`
	Latitude         bool `json:"latitude"`
	Longitude        bool `json:"longitude"`
	InputID          bool `json:"input_id"`
	PopularTimes     bool `json:"popular_times"`
	Cid              bool `json:"cid"`
	Status           bool `json:"status"`
	Descriptions     bool `json:"descriptions"`
	ReviewsLink      bool `json:"reviews_link"`
	Thumbnail        bool `json:"thumbnail"`
	DataID           bool `json:"data_id"`
	StreetViewURL    bool `json:"street_view_url"`
	PlaceID          bool `json:"place_id"`
	Images           bool `json:"images"`
	Reservations     bool `json:"reservations"`
	OrderOnline      bool `json:"order_online"`
	Menu             bool `json:"menu"`
	Owner            bool `json:"owner"`
	CompleteAddress  bool `json:"complete_address"`
}

// DefaultFieldSelection enables every optional field.
func DefaultFieldSelection() FieldSelection {
	return FieldSelection{
		Link:             true,
		Title:            true,
		Category:         true,
		Address:          true,
		OpenHours:        true,
		Website:          true,
		Phone:            true,
		PlusCode:         true,
		ReviewCount:      true,
		ReviewRating:     true,
		ReviewsPerRating: true,
		Latitude:         true,
		Longitude:        true,
		InputID:          true,
		PopularTimes:     true,
		Cid:              true,
		Status:           true,
		Descriptions:     true,
		ReviewsLink:      true,
		Thumbnail:        true,
		DataID:           true,
		StreetViewURL:    true,
		PlaceID:          true,
		Images:           true,
		Reservations:     true,
		OrderOnline:      true,
		Menu:             true,
		Owner:            true,
		CompleteAddress:  true,
	}
}

// ResultFieldSelection returns the selection to use for CSV exports.
func (d JobData) ResultFieldSelection() FieldSelection {
	return d.FieldSelection
}

func (s FieldSelection) includesHeader(header string) bool {
	switch strings.ToLower(header) {
	case "link":
		return s.Link
	case "title":
		return s.Title
	case "category":
		return s.Category
	case "address":
		return s.Address
	case "open_hours":
		return s.OpenHours
	case "website":
		return s.Website
	case "phone":
		return s.Phone
	case "plus_code":
		return s.PlusCode
	case "review_count":
		return s.ReviewCount
	case "review_rating":
		return s.ReviewRating
	case "reviews_per_rating":
		return s.ReviewsPerRating
	case "latitude":
		return s.Latitude
	case "longitude":
		return s.Longitude
	case "input_id":
		return s.InputID
	case "popular_times":
		return s.PopularTimes
	case "cid":
		return s.Cid
	case "status":
		return s.Status
	case "descriptions":
		return s.Descriptions
	case "reviews_link":
		return s.ReviewsLink
	case "thumbnail":
		return s.Thumbnail
	case "data_id":
		return s.DataID
	case "street_view_url":
		return s.StreetViewURL
	case "place_id":
		return s.PlaceID
	case "images":
		return s.Images
	case "reservations":
		return s.Reservations
	case "order_online":
		return s.OrderOnline
	case "menu":
		return s.Menu
	case "owner":
		return s.Owner
	case "complete_address":
		return s.CompleteAddress
	default:
		return true
	}
}
