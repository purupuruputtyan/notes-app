-- ローカル開発用テストユーザー（本番では実行しないこと）
-- ログイン: dev@example.com / Dev1!pass（ユーザー作成APIのパスワードルールと同じ）
INSERT INTO users (id, nick_name, email, password_hash, icon_image, is_active)
VALUES (
    '10000000-0000-4000-8000-000000000001',
    'dev_user',
    'dev@example.com',
    '$2a$10$HSoc8C8rrlG6nPeWv3aibOmpR2n/ws8tSwHd3WLRf3KRdNo8G8rZ6',
    NULL,
    true
)
ON CONFLICT (email) DO NOTHING;

-- ローカル開発用テストノート（上記 dev_user に紐づく 10 件）
INSERT INTO notes (id, user_id, title, content, created_at, updated_at)
VALUES
    (
        '20000000-0000-4000-8000-000000000001',
        '10000000-0000-4000-8000-000000000001',
        'はじめてのメモ',
        'notes-app のローカル開発用シードです。',
        NOW() - INTERVAL '9 days',
        NOW() - INTERVAL '9 days'
    ),
    (
        '20000000-0000-4000-8000-000000000002',
        '10000000-0000-4000-8000-000000000001',
        '買い物リスト',
        '牛乳、卵、パン、コーヒー豆',
        NOW() - INTERVAL '8 days',
        NOW() - INTERVAL '8 days'
    ),
    (
        '20000000-0000-4000-8000-000000000003',
        '10000000-0000-4000-8000-000000000001',
        'API 設計メモ',
        'GET /notes は認証必須。user_id は JWT から取得する。',
        NOW() - INTERVAL '7 days',
        NOW() - INTERVAL '6 days'
    ),
    (
        '20000000-0000-4000-8000-000000000004',
        '10000000-0000-4000-8000-000000000001',
        'リファクタ TODO',
        'ハンドラのエラーログ、ルートへの auth ミドルウェア適用を確認する。',
        NOW() - INTERVAL '6 days',
        NOW() - INTERVAL '5 days'
    ),
    (
        '20000000-0000-4000-8000-000000000005',
        '10000000-0000-4000-8000-000000000001',
        '空の content ではない例',
        'content カラムは NOT NULL のため、空文字ではなく短文を入れる。',
        NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '5 days'
    ),
    (
        '20000000-0000-4000-8000-000000000006',
        '10000000-0000-4000-8000-000000000001',
        '長めのタイトルでも問題ないか確認用のノートタイトルサンプル',
        '本文は普通の長さです。',
        NOW() - INTERVAL '4 days',
        NOW() - INTERVAL '3 days'
    ),
    (
        '20000000-0000-4000-8000-000000000007',
        '10000000-0000-4000-8000-000000000001',
        '週次レビュー',
        '今週: 一覧 API、リポジトリ、ユースケース、ハンドラまで実装。',
        NOW() - INTERVAL '3 days',
        NOW() - INTERVAL '2 days'
    ),
    (
        '20000000-0000-4000-8000-000000000008',
        '10000000-0000-4000-8000-000000000001',
        'curl メモ',
        'curl -H "Authorization: Bearer <token>" http://localhost:8080/notes',
        NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '1 day'
    ),
    (
        '20000000-0000-4000-8000-000000000009',
        '10000000-0000-4000-8000-000000000001',
        'アーカイブ候補',
        '古いメモは将来ソフトデリートやタグ付けを検討。',
        NOW() - INTERVAL '1 day',
        NOW() - INTERVAL '12 hours'
    ),
    (
        '20000000-0000-4000-8000-000000000010',
        '10000000-0000-4000-8000-000000000001',
        '最新のメモ',
        'created_at DESC で並ぶことを確認するための一番新しい行。',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;
