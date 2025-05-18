### Tạo database mysql in docker

```bash

docker run --name payday -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=payday -p 3309:3306 -d mysql

```

### Chạy project

```bash

CompileDaemon -command="./payday"

```

### Chạy Cloud Flared

Download library

```bash

npm i

```

Create quick Tunnel!

```bash

cloudflared tunnel --url http://localhost:8080

```

Lấy tunnel đã tạo được thay thay cho tunnel trong .env
URL_FREE_CLOUD_FLARE="Your Tunnel"

Thay đổi variables trong Postman
url = Your Tunnel

### Tạo mới một giao dịch

Truy cập {{url}}/payment_bank với phương thức POST
{{url}} = [Your Tunnel](#Chạy-Cloud-Flared)

# Body

```json
{
  "amount": 1000000,
  "info": "Test order cho toi"
}
```

Nhận payment_url và mở trên trình duyệt và hoàn tất thanh toán

### Kiểm tra trạng thái giao dịch

Truy {{url}}/check_order/{{transaction_id}} với phương thức GET
{{url}} = [Your Tunnel](#Chạy-Cloud-Flared)
{{transaction_id}} = Mã giao dịch. VD: AP251458220631

### Hoàn tiền

Truy cập {{url}}/refund với phương thức POST
{{url}} = [Your Tunnel](#Chạy-Cloud-Flared)

# Body

```json
{
  "transactionId": "{{transaction_id}}",
  "Reason": "Het tien roi tra lai de"
}
```

{{transaction_id}} = Mã giao dịch. VD: AP251458220631
