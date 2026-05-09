package marketsurge

import "context"

// queryNavTree is the GraphQL query for the NavTree endpoint.
const queryNavTree = `query NavTree($site: Site!, $treeType: NavTreeTypeInput!) {
  user {
    navTree(site: $site, treeType: $treeType) {
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

// ---------------------------------------------------------------------------
// Defaults
// ---------------------------------------------------------------------------

const (
	// DefaultNavTreeSite is the default site for NavTree queries.
	DefaultNavTreeSite = "marketsurge"
	// DefaultNavTreeTreeType is the default tree type for NavTree queries.
	DefaultNavTreeTreeType = "MSR_NAV"
)

// ---------------------------------------------------------------------------
// NavTree request types
// ---------------------------------------------------------------------------

// NavTreeRequest holds parameters for the NavTree query.
type NavTreeRequest struct {
	Site     string
	TreeType string
}

// NewNavTreeRequest creates a NavTreeRequest with sensible defaults.
func NewNavTreeRequest() NavTreeRequest {
	return NavTreeRequest{
		Site:     DefaultNavTreeSite,
		TreeType: DefaultNavTreeTreeType,
	}
}

// navTreeVariables holds the GraphQL variables sent with the query.
type navTreeVariables struct {
	Site     string `json:"site"`
	TreeType string `json:"treeType"`
}

// ---------------------------------------------------------------------------
// NavTree response types
// ---------------------------------------------------------------------------

// NavTreeResponse is the top-level response from the NavTree query.
type NavTreeResponse struct {
	User *NavTreeUser `json:"user"`
}

// NavTreeUser wraps the user-scoped NavTree result.
type NavTreeUser struct {
	NavTree []NavTreeNode `json:"navTree"`
}

// NavTreeNode represents a single node in the navigation tree, which may be
// either a folder or a leaf. Folder-only fields (Children, ContentType) are
// nil for leaf nodes; leaf-only fields (URL, ReferenceID) are nil for folders.
type NavTreeNode struct {
	ID          *string            `json:"id"`
	Name        *string            `json:"name"`
	ParentID    *string            `json:"parentId"`
	Type        *string            `json:"type"`
	Children    []NavTreeChildNode `json:"children"`
	ContentType *string            `json:"contentType"`
	TreeType    *string            `json:"treeType"`
	URL         *string            `json:"url"`
	ReferenceID *string            `json:"referenceId"`
}

// NavTreeChildNode represents a child node summary within a folder.
type NavTreeChildNode struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
	Type *string `json:"type"`
}

// ---------------------------------------------------------------------------
// Client methods
// ---------------------------------------------------------------------------

// NavTree fetches the full navigation tree for the authenticated user,
// including watchlists, screens, and built-in reports.
func (c *Client) NavTree(ctx context.Context, req NavTreeRequest) (*NavTreeResponse, error) {
	vars := navTreeVariables(req)

	var resp NavTreeResponse
	if err := c.doGraphQL(ctx, "NavTree", vars, queryNavTree, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
