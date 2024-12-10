package models

import (
	"encoding/json"
	"time"

	"github.com/mrrizkin/omniscan/pkg/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"gorm.io/gorm"
)

type EStatementMetadata struct {
	ID                 uint           `json:"id"            gorm:"primary_key"`
	CreatedAt          *time.Time     `json:"created_at"`
	UpdatedAt          *time.Time     `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at"    gorm:"index"`
	EStatementID       uint           `json:"e_statement_id" gorm:"index"`
	FileName           string         `json:"source,omitempty"`
	Version            string         `json:"version"`
	PageCount          int            `json:"pageCount"`
	Boundaries         string         `json:"pageBoundaries,omitempty"`
	Dimensions         string         `json:"pageSizes,omitempty"`
	Title              string         `json:"title"`
	Author             string         `json:"author"`
	Subject            string         `json:"subject"`
	Producer           string         `json:"producer"`
	Creator            string         `json:"creator"`
	CreationDate       string         `json:"creationDate"`
	ModificationDate   string         `json:"modificationDate"`
	PageMode           string         `json:"pageMode,omitempty"`
	PageLayout         string         `json:"pageLayout,omitempty"`
	Keywords           string         `json:"keywords"`
	Properties         string         `json:"properties"`
	Tagged             bool           `json:"tagged"`
	Hybrid             bool           `json:"hybrid"`
	Linearized         bool           `json:"linearized"`
	UsingXRefStreams   bool           `json:"usingXRefStreams"`
	UsingObjectStreams bool           `json:"usingObjectStreams"`
	Watermarked        bool           `json:"watermarked"`
	Thumbnails         bool           `json:"thumbnails"`
	Form               bool           `json:"form"`
	Signatures         bool           `json:"signatures"`
	AppendOnly         bool           `json:"appendOnly"`
	Outlines           bool           `json:"bookmarks"`
	Names              bool           `json:"names"`
	Encrypted          bool           `json:"encrypted"`
	Permissions        int            `json:"permissions"`
	Attachments        string         `json:"attachments,omitempty"`
	UnitString         string         `json:"unit"`
	XMLMetadata        string         `json:"xmlmetadata,omitempty"`
}

func IntoEstatementMetadata(eStatementID uint, meta *pdf.Metadata) (*EStatementMetadata, error) {
	boundaries, err := json.Marshal(meta.Boundaries)
	if err != nil {
		return nil, err
	}

	dimensions, err := json.Marshal(meta.Dimensions)
	if err != nil {
		return nil, err
	}

	modificationDate, err := json.Marshal(meta.ModificationDate)
	if err != nil {
		return nil, err
	}

	keywords, err := json.Marshal(meta.Keywords)
	if err != nil {
		return nil, err
	}

	properties, err := json.Marshal(meta.Properties)
	if err != nil {
		return nil, err
	}

	attachments, err := json.Marshal(meta.Attachments)
	if err != nil {
		return nil, err
	}

	xmlMetadata, err := json.Marshal(meta.XMLMetadata)
	if err != nil {
		return nil, err
	}

	return &EStatementMetadata{
		EStatementID:       eStatementID,
		FileName:           meta.FileName,
		Version:            meta.Version,
		PageCount:          meta.PageCount,
		Boundaries:         string(boundaries),
		Dimensions:         string(dimensions),
		Title:              meta.Title,
		Author:             meta.Author,
		Subject:            meta.Subject,
		Producer:           meta.Producer,
		Creator:            meta.Creator,
		CreationDate:       meta.CreationDate,
		ModificationDate:   string(modificationDate),
		PageMode:           meta.PageMode,
		PageLayout:         meta.PageLayout,
		Keywords:           string(keywords),
		Properties:         string(properties),
		Tagged:             meta.Tagged,
		Hybrid:             meta.Hybrid,
		Linearized:         meta.Linearized,
		UsingXRefStreams:   meta.UsingXRefStreams,
		UsingObjectStreams: meta.UsingObjectStreams,
		Watermarked:        meta.Watermarked,
		Thumbnails:         meta.Thumbnails,
		Form:               meta.Form,
		Signatures:         meta.Signatures,
		AppendOnly:         meta.AppendOnly,
		Outlines:           meta.Outlines,
		Names:              meta.Names,
		Encrypted:          meta.Encrypted,
		Permissions:        meta.Permissions,
		Attachments:        string(attachments),
		UnitString:         meta.UnitString,
		XMLMetadata:        string(xmlMetadata),
	}, nil
}

func (m *EStatementMetadata) ToMetadata() (*pdf.Metadata, error) {
	boundaries := make(map[string]model.PageBoundaries)
	if err := json.Unmarshal([]byte(m.Boundaries), &boundaries); err != nil {
		return nil, err
	}

	dimensions := []types.Dim{}
	if err := json.Unmarshal([]byte(m.Dimensions), &dimensions); err != nil {
		return nil, err
	}

	keywords := []string{}
	if err := json.Unmarshal([]byte(m.Keywords), &keywords); err != nil {
		return nil, err
	}

	properties := make(map[string]string)
	if err := json.Unmarshal([]byte(m.Properties), &properties); err != nil {
		return nil, err
	}

	attachments := []model.Attachment{}
	if err := json.Unmarshal([]byte(m.Attachments), &attachments); err != nil {
		return nil, err
	}

	pdfInfo := &pdfcpu.PDFInfo{
		FileName:           m.FileName,
		Version:            m.Version,
		PageCount:          m.PageCount,
		Boundaries:         boundaries,
		Dimensions:         dimensions,
		Title:              m.Title,
		Author:             m.Author,
		Subject:            m.Subject,
		Producer:           m.Producer,
		Creator:            m.Creator,
		CreationDate:       m.CreationDate,
		ModificationDate:   m.ModificationDate,
		PageMode:           m.PageMode,
		PageLayout:         m.PageLayout,
		Keywords:           keywords,
		Properties:         properties,
		Tagged:             m.Tagged,
		Hybrid:             m.Hybrid,
		Linearized:         m.Linearized,
		UsingXRefStreams:   m.UsingXRefStreams,
		UsingObjectStreams: m.UsingObjectStreams,
		Watermarked:        m.Watermarked,
		Thumbnails:         m.Thumbnails,
		Form:               m.Form,
		Signatures:         m.Signatures,
		AppendOnly:         m.AppendOnly,
		Outlines:           m.Outlines,
		Names:              m.Names,
		Encrypted:          m.Encrypted,
		Permissions:        m.Permissions,
		Attachments:        attachments,
		UnitString:         m.UnitString,
	}

	return &pdf.Metadata{
		PDFInfo: pdfInfo,
	}, nil
}
