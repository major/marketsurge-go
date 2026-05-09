package marketsurge

import "context"

// queryCoachTree is the GraphQL query for the CoachTree endpoint.
const queryCoachTree = `query CoachTree($site: Site!, $treeType: NavTreeTypeInput!) {
  user {
    watchlists: coachTree(
      coachTreeType: WATCHLIST
      site: $site
      treeType: $treeType
    ) {
      ... on NavTreeFolder {
        id
        name
        parentId
        type
        children {
          ... on NavTreeFolder { id name type }
          ... on NavTreeLeaf { id name type }
        }
        contentType
        treeType
      }
      ... on NavTreeLeaf {
        id
        name
        parentId
        type
        url
        treeType
        referenceId
      }
    }
    screens: coachTree(
      coachTreeType: SCREEN
      site: $site
      treeType: $treeType
    ) {
      ... on NavTreeFolder {
        id
        name
        parentId
        type
        children {
          ... on NavTreeFolder { id name type }
          ... on NavTreeLeaf { id name type }
        }
        contentType
        treeType
      }
      ... on NavTreeLeaf {
        id
        name
        parentId
        type
        url
        treeType
        referenceId
      }
    }
  }
}`

// queryIndustryGroupRS is the GraphQL query for the IndustryGroupRS endpoint.
const queryIndustryGroupRS = `query IndustryGroupRS($symbols: [String!]!, $symbolDialectType: MDSymbolDialectType!) {
  marketData(symbols: $symbols, symbolDialectType: $symbolDialectType) {
    originRequest {
      symbol
    }
    industry {
      groupRS(where: { periodOffset: { eq: CURRENT }, period: { eq: P6M } }) {
        value
      }
    }
  }
}`

// ---------------------------------------------------------------------------
// Defaults
// ---------------------------------------------------------------------------

const (
	// DefaultCoachTreeSite is the default site for CoachTree queries.
	DefaultCoachTreeSite = "marketsurge"
	// DefaultCoachTreeTreeType is the default tree type for CoachTree queries.
	DefaultCoachTreeTreeType = "MSR_NAV"
	// DefaultIndustryGroupRSSymbolDialectType is the default symbol dialect
	// for IndustryGroupRS queries.
	DefaultIndustryGroupRSSymbolDialectType = "CHARTING"
)

// ---------------------------------------------------------------------------
// CoachTree request types
// ---------------------------------------------------------------------------

// CoachTreeRequest holds parameters for the CoachTree query.
type CoachTreeRequest struct {
	Site     string
	TreeType string
}

// NewCoachTreeRequest creates a CoachTreeRequest with sensible defaults.
func NewCoachTreeRequest() CoachTreeRequest {
	return CoachTreeRequest{
		Site:     DefaultCoachTreeSite,
		TreeType: DefaultCoachTreeTreeType,
	}
}

// coachTreeVariables holds the GraphQL variables sent with the query.
type coachTreeVariables struct {
	Site     string `json:"site"`
	TreeType string `json:"treeType"`
}

// ---------------------------------------------------------------------------
// IndustryGroupRS request types
// ---------------------------------------------------------------------------

// IndustryGroupRSRequest holds parameters for the IndustryGroupRS query.
type IndustryGroupRSRequest struct {
	Symbols           []string
	SymbolDialectType string
}

// NewIndustryGroupRSRequest creates an IndustryGroupRSRequest with sensible
// defaults for the given symbols.
func NewIndustryGroupRSRequest(symbols ...string) IndustryGroupRSRequest {
	return IndustryGroupRSRequest{
		Symbols:           symbols,
		SymbolDialectType: DefaultIndustryGroupRSSymbolDialectType,
	}
}

// industryGroupRSVariables holds the GraphQL variables sent with the query.
type industryGroupRSVariables struct {
	Symbols           []string `json:"symbols"`
	SymbolDialectType string   `json:"symbolDialectType"`
}

// ---------------------------------------------------------------------------
// CoachTree response types
// ---------------------------------------------------------------------------

// CoachTreeResponse is the top-level response from the CoachTree query.
type CoachTreeResponse struct {
	User *CoachTreeUser `json:"user"`
}

// CoachTreeUser wraps the user-scoped CoachTree result.
type CoachTreeUser struct {
	Watchlists []CoachTreeNode `json:"watchlists"`
	Screens    []CoachTreeNode `json:"screens"`
}

// CoachTreeNode represents a single node in the coach tree, which may be
// either a folder or a leaf. Folder-only fields (Children, ContentType) are
// nil for leaf nodes; leaf-only fields (URL, ReferenceID) are nil for folders.
type CoachTreeNode struct {
	ID          *string              `json:"id"`
	Name        *string              `json:"name"`
	ParentID    *string              `json:"parentId"`
	Type        *string              `json:"type"`
	Children    []CoachTreeChildNode `json:"children"`
	ContentType *string              `json:"contentType"`
	TreeType    *string              `json:"treeType"`
	URL         *string              `json:"url"`
	ReferenceID *string              `json:"referenceId"`
}

// CoachTreeChildNode represents a child node summary within a folder.
type CoachTreeChildNode struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
	Type *string `json:"type"`
}

// ---------------------------------------------------------------------------
// IndustryGroupRS response types
// ---------------------------------------------------------------------------

// IndustryGroupRSResponse is the top-level response from the IndustryGroupRS query.
type IndustryGroupRSResponse struct {
	MarketData []IndustryGroupRSItem `json:"marketData"`
}

// IndustryGroupRSItem represents a single symbol's industry group RS data.
type IndustryGroupRSItem struct {
	OriginRequest *IndustryGroupRSOriginRequest `json:"originRequest"`
	Industry      *IndustryGroupRSIndustry      `json:"industry"`
}

// IndustryGroupRSOriginRequest holds the original request symbol.
type IndustryGroupRSOriginRequest struct {
	Symbol *string `json:"symbol"`
}

// IndustryGroupRSIndustry holds industry data including group RS values.
type IndustryGroupRSIndustry struct {
	GroupRS []IndustryGroupRSGroupRS `json:"groupRS"`
}

// IndustryGroupRSGroupRS holds a single group RS rating value.
type IndustryGroupRSGroupRS struct {
	Value *int `json:"value"`
}

// ---------------------------------------------------------------------------
// Client methods
// ---------------------------------------------------------------------------

// CoachTree fetches the navigation tree for watchlists and screens.
func (c *Client) CoachTree(ctx context.Context, req CoachTreeRequest) (*CoachTreeResponse, error) {
	vars := coachTreeVariables(req)

	var resp CoachTreeResponse
	if err := c.doGraphQL(ctx, "CoachTree", vars, queryCoachTree, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IndustryGroupRS fetches industry group relative strength data for the
// requested symbols.
func (c *Client) IndustryGroupRS(ctx context.Context, req IndustryGroupRSRequest) (*IndustryGroupRSResponse, error) {
	vars := industryGroupRSVariables(req)

	var resp IndustryGroupRSResponse
	if err := c.doGraphQL(ctx, "IndustryGroupRS", vars, queryIndustryGroupRS, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
