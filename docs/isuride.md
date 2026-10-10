classDiagram
direction BT
class node9 {
   varchar(26) chair_id  /* 椅子ID */
   int latitude  /* 経度 */
   int longitude  /* 緯度 */
   datetime(6) created_at  /* 登録日時 */
   varchar(26) id
}
class node6 {
   int speed  /* 移動速度 */
   varchar(50) name  /* 椅子モデル名 */
}
class node8 {
   varchar(26) owner_id  /* オーナーID */
   varchar(30) name  /* 椅子の名前 */
   text model  /* 椅子のモデル */
   tinyint(1) is_active  /* 配椅子受付中かどうか */
   varchar(255) access_token  /* アクセストークン */
   datetime(6) created_at  /* 登録日時 */
   datetime(6) updated_at  /* 更新日時 */
   varchar(26) id  /* 椅子ID */
}
class node0 {
   int discount  /* 割引額 */
   datetime(6) created_at  /* 付与日時 */
   varchar(26) used_by  /* クーポンが適用されたライドのID */
   varchar(26) user_id  /* 所有しているユーザーのID */
   varchar(255) code  /* クーポンコード */
}
class node3 {
   varchar(30) name  /* オーナー名 */
   varchar(255) access_token  /* アクセストークン */
   varchar(255) chair_register_token  /* 椅子登録トークン */
   datetime(6) created_at  /* 登録日時 */
   datetime(6) updated_at  /* 更新日時 */
   varchar(26) id  /* オーナーID */
}
class node2 {
   varchar(255) token  /* 決済トークン */
   datetime(6) created_at  /* 登録日時 */
   varchar(26) user_id  /* ユーザーID */
}
class node4 {
   varchar(26) ride_id  /* ライドID */
   enum('matching', 'enroute', 'pickup', 'carrying', 'arrived', 'completed') status  /* 状態 */
   datetime(6) created_at  /* 状態変更日時 */
   datetime(6) app_sent_at  /* ユーザーへの状態通知日時 */
   datetime(6) chair_sent_at  /* 椅子への状態通知日時 */
   varchar(26) id
}
class node7 {
   varchar(26) user_id  /* ユーザーID */
   varchar(26) chair_id  /* 割り当てられた椅子ID */
   pickup_latitude  /* 配車位置(経度) */ int
   pickup_longitude  /* 配車位置(緯度) */ int
   destination_latitude  /* 目的地(経度) */ int
   destination_longitude  /* 目的地(緯度) */ int
   int evaluation  /* 評価 */
   datetime(6) created_at  /* 要求日時 */
   datetime(6) updated_at  /* 状態更新日時 */
   varchar(26) id  /* ライドID */
}
class node1 {
   text value  /* 設定値 */
   varchar(30) name  /* 設定名 */
}
class node5 {
   varchar(30) username  /* ユーザー名 */
   firstname  /* 本名(名前) */ varchar(30)
   lastname  /* 本名(名字) */ varchar(30)
   varchar(30) date_of_birth  /* 生年月日 */
   varchar(255) access_token  /* アクセストークン */
   varchar(30) invitation_code  /* 招待トークン */
   datetime(6) created_at  /* 登録日時 */
   datetime(6) updated_at  /* 更新日時 */
   varchar(26) id  /* ユーザーID */
}

node9  -->  node8 : chair_id:id
node8  -->  node3 : owner_id:id
node0  -->  node5 : user_id:id
node2  -->  node5 : user_id:id
node4  -->  node7 : ride_id:id
node7  -->  node8 : chair_id:id
node7  -->  node5 : user_id:id
