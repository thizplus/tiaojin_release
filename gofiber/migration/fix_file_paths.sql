-- แก้ไข full URL → path only
-- ตัด prefix: https://tiaojin.s3.ap-southeast-1.idrivee2.com/

-- 1. Company Profiles
UPDATE company_profiles SET 
    logo = REPLACE(logo, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', ''),
    atta = REPLACE(atta, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', ''),
    dbd = REPLACE(dbd, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', ''),
    pp20 = REPLACE(pp20, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', '')
WHERE logo LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%'
   OR atta LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%'
   OR dbd LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%'
   OR pp20 LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%';

-- 2. Banners
UPDATE banners SET 
    img = REPLACE(img, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', '')
WHERE img LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%';

-- 3. Cities
UPDATE cities SET 
    img = REPLACE(img, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', '')
WHERE img LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%';

-- 4. Tours
UPDATE tours SET 
    thumbnail = REPLACE(thumbnail, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', ''),
    file_pdf = REPLACE(file_pdf, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', ''),
    file_word = REPLACE(file_word, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', '')
WHERE thumbnail LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%'
   OR file_pdf LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%'
   OR file_word LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%';

-- 5. Periods (JSONB - ต้องใช้ regexp_replace)
UPDATE periods SET 
    files = REPLACE(files::text, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', '')::jsonb
WHERE files::text LIKE '%https://tiaojin.s3.ap-southeast-1.idrivee2.com/%';

-- 6. Payment Slips
UPDATE payment_slips SET 
    image = REPLACE(image, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', '')
WHERE image LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%';

-- 7. Passenger Passports
UPDATE passenger_passports SET 
    image = REPLACE(image, 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/', '')
WHERE image LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2.com/%';

-- ตรวจสอบผลลัพธ์
SELECT 'company_profiles' as tbl, COUNT(*) as remaining FROM company_profiles WHERE logo LIKE 'https://%' OR atta LIKE 'https://%'
UNION ALL
SELECT 'banners', COUNT(*) FROM banners WHERE img LIKE 'https://tiaojin.s3%'
UNION ALL
SELECT 'cities', COUNT(*) FROM cities WHERE img LIKE 'https://tiaojin.s3%'
UNION ALL
SELECT 'tours', COUNT(*) FROM tours WHERE thumbnail LIKE 'https://tiaojin.s3%';
