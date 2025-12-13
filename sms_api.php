<?php
class ApiClient
{
  private $apiKey;
  private $apiSecret;
  private $url;

  public function __construct()
  {
    // Hardcode credentials & full API URL here
    $this->apiKey = 'client123SWDsms';
    $this->apiSecret = '2946fe9f50e4d5cc7f6e7fec1779d1eb032d6c2c18b2b9f31fd87a60706258a5';
    $this->url = 'https://console.fiberathome.net/sms/';
  }

  public function sendRequest($payload = [])
  {
    
    $timestamp = time();
    $body = json_encode($payload);
    $signature = hash_hmac('sha256', $body . $timestamp, $this->apiSecret);

    $ch = curl_init();
    curl_setopt_array($ch, [
      CURLOPT_URL => $this->url,
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_POST => true,
      CURLOPT_HTTPHEADER => [
        "Content-Type: application/json",
        "X-API-KEY: {$this->apiKey}",
        "X-TIMESTAMP: {$timestamp}",
        "X-SIGNATURE: {$signature}"
      ],
      CURLOPT_POSTFIELDS => $body,
    ]);

    $response = curl_exec($ch);
    if (curl_errno($ch)) {
      $error = curl_error($ch);
      curl_close($ch);
      return ['error' => $error];
    }

    curl_close($ch);
    $decoded = json_decode($response, true);
    return $decoded ?: $response;
  }
}
