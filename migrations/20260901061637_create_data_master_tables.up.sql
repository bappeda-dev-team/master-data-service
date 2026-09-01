CREATE TABLE data_master_sync (
    id CHAR(36) NOT NULL,

    sync_from VARCHAR(255) NOT NULL,
    sync_type VARCHAR(50) NOT NULL,

    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP NULL,

    status VARCHAR(20) NOT NULL,

    total_data INT NOT NULL DEFAULT 0,
    total_success INT NOT NULL DEFAULT 0,
    total_failed INT NOT NULL DEFAULT 0,

    PRIMARY KEY (id),

    INDEX idx_data_master_sync_source (
        sync_from,
        sync_type,
        started_at
    )
);


CREATE TABLE data_master_opd (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,

    opd_kode VARCHAR(50) NOT NULL,
    opd_level SMALLINT NOT NULL,
    opd_nama VARCHAR(255) NOT NULL,
    opd_jabatan VARCHAR(255) NULL,
    eselon_kode VARCHAR(20) NULL,
    eselon_nama VARCHAR(50) NULL,

    sync_id CHAR(36) NULL,
    sync_from VARCHAR(255) NULL,
    sync_at TIMESTAMP NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    CONSTRAINT uk_data_master_opd_kode
        UNIQUE (opd_kode),

    CONSTRAINT fk_data_master_opd_sync
        FOREIGN KEY (sync_id)
        REFERENCES data_master_sync (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);


CREATE TABLE data_master_jabatan (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,

    opd_kode VARCHAR(50) NOT NULL,

    jenis_jabatan_kode VARCHAR(20) NOT NULL,
    jenis_jabatan_nama VARCHAR(100) NULL,

    jabatan_kode VARCHAR(50) NOT NULL,
    jabatan_nama VARCHAR(255) NOT NULL,

    jenjang_jabatan_kode VARCHAR(20) NULL,
    jenjang_jabatan_nama VARCHAR(100) NULL,

    sync_id CHAR(36) NULL,
    sync_from VARCHAR(255) NULL,
    sync_at TIMESTAMP NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    CONSTRAINT uk_data_master_jabatan
        UNIQUE (
            opd_kode,
            jenis_jabatan_kode,
            jabatan_kode
        ),

    CONSTRAINT fk_data_master_jabatan_opd
        FOREIGN KEY (opd_kode)
        REFERENCES data_master_opd (opd_kode)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_data_master_jabatan_sync
        FOREIGN KEY (sync_id)
        REFERENCES data_master_sync (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL,

    INDEX idx_data_master_jabatan_opd (opd_kode)
);


CREATE TABLE data_master_pegawai (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,

    pegawai_nip VARCHAR(30) NOT NULL,

    pegawai_nama VARCHAR(255) NOT NULL,

    pegawai_jenis_kelamin VARCHAR(30) NULL,

    pegawai_lahir_tempat VARCHAR(100) NULL,
    pegawai_lahir_tanggal DATE NULL,

    pegawai_pensiun_tmt DATE NULL,

    pegawai_aktif_id VARCHAR(20) NULL,
    pegawai_status VARCHAR(100) NULL,
    pegawai_jenis VARCHAR(255) NULL,

    pegawai_golongan_kode VARCHAR(20) NULL,
    pegawai_golongan_nama VARCHAR(50) NULL,
    pegawai_golongan_pangkat VARCHAR(100) NULL,
    pegawai_golongan_tmt DATE NULL,

    pegawai_jenis_jabatan_kode VARCHAR(20) NULL,
    pegawai_jenis_jabatan_nama VARCHAR(100) NULL,

    pegawai_jabatan_kode VARCHAR(50) NULL,
    pegawai_jabatan_terakhir VARCHAR(255) NULL,

    pegawai_eselon_kode VARCHAR(20) NULL,
    pegawai_eselon_nama VARCHAR(50) NULL,

    opd_kode VARCHAR(50) NOT NULL,

    opd_level SMALLINT NULL,
    opd_nama VARCHAR(255) NULL,

    sync_id CHAR(36) NULL,
    sync_from VARCHAR(255) NULL,
    sync_at TIMESTAMP NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    CONSTRAINT uk_data_master_pegawai_nip
        UNIQUE (pegawai_nip),

    CONSTRAINT fk_data_master_pegawai_opd
        FOREIGN KEY (opd_kode)
        REFERENCES data_master_opd (opd_kode)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_data_master_pegawai_sync
        FOREIGN KEY (sync_id)
        REFERENCES data_master_sync (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL,

    INDEX idx_data_master_pegawai_opd (opd_kode),

    INDEX idx_data_master_pegawai_jabatan (
        opd_kode,
        pegawai_jenis_jabatan_kode,
        pegawai_jabatan_kode
    )
);
