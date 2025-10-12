CREATE TABLE IF NOT EXISTS RepoData (
    id double precision,
    name VARCHAR(100),
    user_kind VARCHAR(100),
    user_name VARCHAR(100),
    is_private bool,
    is_fork bool
)
