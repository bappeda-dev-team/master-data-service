CREATE TABLE data_master_sync (
    id UUID PRIMARY KEY,

    sync_from VARCHAR(255) NOT NULL,
    sync_type VARCHAR(50) NOT NULL,

    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP NULL,

    status VARCHAR(20) NOT NULL,

    total_data INTEGER NOT NULL DEFAULT 0,
    total_success INTEGER NOT NULL DEFAULT 0,
    total_failed INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_data_master_sync_source
    ON data_master_sync (
        sync_from,
        sync_type,
        started_at
    );


CREATE TABLE data_master_opd (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    opd_kode VARCHAR(50) NOT NULL,
    opd_level SMALLINT NOT NULL,
    opd_nama VARCHAR(255) NOT NULL,
    opd_jabatan VARCHAR(255),
    eselon_kode VARCHAR(20),
    eselon_nama VARCHAR(50),

    sync_id UUID,
    sync_from VARCHAR(255),
    sync_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_data_master_opd
        UNIQUE (sync_id, opd_kode),

    CONSTRAINT fk_data_master_opd_sync
        FOREIGN KEY (sync_id)
        REFERENCES data_master_sync (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

CREATE INDEX idx_data_master_opd_kode
    ON data_master_opd (opd_kode);

CREATE INDEX idx_data_master_opd_sync_kode
    ON data_master_opd (sync_id, opd_kode);


CREATE TABLE data_master_jabatan (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    opd_kode VARCHAR(50) NOT NULL,

    jenis_jabatan_kode VARCHAR(20) NOT NULL,
    jenis_jabatan_nama VARCHAR(100),

    jabatan_kode VARCHAR(50) NOT NULL,
    jabatan_nama VARCHAR(255) NOT NULL,

    jenjang_jabatan_kode VARCHAR(20),
    jenjang_jabatan_nama VARCHAR(100),

    sync_id UUID,
    sync_from VARCHAR(255),
    sync_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_data_master_jabatan
        UNIQUE (
            sync_id,
            opd_kode,
            jenis_jabatan_kode,
            jabatan_kode
        ),

    CONSTRAINT fk_data_master_jabatan_sync
        FOREIGN KEY (sync_id)
        REFERENCES data_master_sync (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

CREATE INDEX idx_data_master_jabatan_opd
    ON data_master_jabatan (opd_kode);

CREATE INDEX idx_data_master_jabatan_sync_opd
    ON data_master_jabatan (sync_id, opd_kode);


CREATE TABLE data_master_pegawai (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    opd_kode VARCHAR(50) NOT NULL,

    pegawai_nip VARCHAR(30) NOT NULL,

    pegawai_nama VARCHAR(255) NOT NULL,

    pegawai_jenis_kelamin VARCHAR(30),

    pegawai_lahir_tempat VARCHAR(100),
    pegawai_lahir_tanggal DATE,

    pegawai_pensiun_tmt DATE,

    pegawai_aktif_id VARCHAR(20),
    pegawai_status VARCHAR(100),
    pegawai_jenis VARCHAR(255),

    pegawai_golongan_kode VARCHAR(20),
    pegawai_golongan_nama VARCHAR(50),
    pegawai_golongan_pangkat VARCHAR(100),
    pegawai_golongan_tmt DATE,

    pegawai_jenis_jabatan_kode VARCHAR(20),
    pegawai_jenis_jabatan_nama VARCHAR(100),

    pegawai_jabatan_kode VARCHAR(50),
    pegawai_jabatan_terakhir VARCHAR(255),

    pegawai_eselon_kode VARCHAR(20),
    pegawai_eselon_nama VARCHAR(50),

    opd_level SMALLINT,
    opd_nama VARCHAR(255),

    sync_id UUID,
    sync_from VARCHAR(255),
    sync_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_data_master_pegawai
        UNIQUE (sync_id, pegawai_nip),

    CONSTRAINT fk_data_master_pegawai_sync
        FOREIGN KEY (sync_id)
        REFERENCES data_master_sync (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

CREATE INDEX idx_data_master_pegawai_opd
    ON data_master_pegawai (opd_kode);

CREATE INDEX idx_data_master_pegawai_sync_opd
    ON data_master_pegawai (sync_id, opd_kode);

CREATE INDEX idx_data_master_pegawai_jabatan
    ON data_master_pegawai (
        sync_id,
        opd_kode,
        pegawai_jenis_jabatan_kode,
        pegawai_jabatan_kode
    );
