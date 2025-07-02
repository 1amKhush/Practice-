
CREATE TABLE IF NOT EXISTS peers (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			peer_id VARCHAR(64) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			ip_address INET,
			is_online BOOLEAN DEFAULT false,
			last_seen TIMESTAMP DEFAULT NOW(),
			created_at TIMESTAMP DEFAULT NOW()
		);

	
		CREATE TABLE IF NOT EXISTS files (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			file_hash VARCHAR(128) UNIQUE NOT NULL,
			filename VARCHAR(512) NOT NULL,
			file_size BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);

		
		CREATE TABLE IF NOT EXISTS peer_files (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			peer_id UUID REFERENCES peers(id) ON DELETE CASCADE,
			file_id UUID REFERENCES files(id) ON DELETE CASCADE,
			announced_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(peer_id, file_id)
		);

		
		CREATE TABLE IF NOT EXISTS trust_scores (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			peer_id UUID UNIQUE REFERENCES peers(id) ON DELETE CASCADE,
			score DECIMAL(3,2) DEFAULT 0.50 CHECK (score >= 0.00 AND score <= 1.00),
			successful_transfers INTEGER DEFAULT 0,
			failed_transfers INTEGER DEFAULT 0,
			updated_at TIMESTAMP DEFAULT NOW()
		);

		
		CREATE TABLE IF NOT EXISTS active_connections (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			requester_id UUID REFERENCES peers(id) ON DELETE CASCADE,
			provider_id UUID REFERENCES peers(id) ON DELETE CASCADE,
			file_id UUID REFERENCES files(id) ON DELETE CASCADE,
			status VARCHAR(20) DEFAULT 'connecting'
				CHECK (status IN ('connecting', 'active', 'completed', 'failed')),
			started_at TIMESTAMP DEFAULT NOW(),
			completed_at TIMESTAMP
		);

		