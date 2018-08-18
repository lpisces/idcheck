<?php
require("./vendor/autoload.php");

$key = "";
$secret = "";
$url = "https://test-idcheck.nle-tech.com/";

$name = "王锐";
$number = "230403198902130446";
$expire = time() + 5;

$orig = $secret . "expire=".$expire."&key=".$key."&name=".urlencode($name)."&number=".$number;
echo $orig;
echo "\n";
$sign = md5($orig);
echo $sign;
echo "\n";

$client = new GuzzleHttp\Client();
$res = $client->request('GET', $url . '/id_check', 
	[ 
		"query" => [
			"name" => $name,
			"number" => $number,
			"expire" => $expire,
			"key" => $key,
			"sign" => $sign,
		],
	]
);

echo $res->getBody();
echo "\n";
