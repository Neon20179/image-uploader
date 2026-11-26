CREATE TABLE embedding_blacklist (
  id UUID PRIMARY KEY,
  vector vector(512) NOT NULL CHECK (vector_norm(vector) > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
