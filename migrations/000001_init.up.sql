CREATE EXTENSION vector;

CREATE TYPE upload_status AS
ENUM('pending', 'uploaded', 'processing', 'ready', 'blocked', 'failed');

CREATE TYPE process_status AS
ENUM('pending', 'completed', 'failed');

CREATE TABLE images (
  id UUID PRIMARY KEY,

  bucket VARCHAR(63) NOT NULL,
  object_key VARCHAR(1024) NOT NULL,

  width INTEGER NOT NULL,
  height INTEGER NOT NULL,
  size INTEGER NOT NULL,
  mime_type VARCHAR(255) NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE uploads (
  id UUID PRIMARY KEY,

  image_id UUID REFERENCES images(id) ON DELETE CASCADE,

  status upload_status NOT NULL DEFAULT 'pending',

  bucket VARCHAR(63) NOT NULL,
  object_key VARCHAR(1024),
  raw_object_key VARCHAR(1024) NOT NULL,

  normalization_status process_status NOT NULL DEFAULT 'pending',
  embedding_status process_status NOT NULL DEFAULT 'pending',
  ocr_status process_status NOT NULL DEFAULT 'pending',

  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE normalization_results (
  upload_id UUID PRIMARY KEY REFERENCES uploads(id) ON DELETE CASCADE,

  width INTEGER NOT NULL,
  height INTEGER NOT NULL,
  size INTEGER NOT NULL,
  mime_type VARCHAR(255) NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE embeddings (
  upload_id UUID PRIMARY KEY REFERENCES uploads(id) ON DELETE CASCADE,

  vector vector(512) NOT NULL,

  block_detected BOOLEAN NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ocr_results (
  upload_id UUID PRIMARY KEY REFERENCES uploads(id) ON DELETE CASCADE,

  text TEXT,

  block_detected BOOLEAN NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
