ALTER TABLE kycs
ADD COLUMN front_photo TEXT,
ADD COLUMN back_photo TEXT;

UPDATE kycs
SET
    front_photo = kyc_data.front_photo,
    back_photo = kyc_data.back_photo
FROM
    kyc_data
WHERE
    kycs.user_id = kyc_data.user_id;
