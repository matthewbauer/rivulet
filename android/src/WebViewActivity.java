package com.bauer.rivulet;
 
import android.app.Activity;
import android.os.Bundle;
import android.webkit.WebView;
 
public class WebViewActivity extends Activity {
 
  private WebView webView;
 
  public void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);
    setContentView(R.layout.webview);
 
    webView = (WebView) findViewById(R.id.webview);
    webView.getSettings().setJavaScriptEnabled(true);
    webView.loadUrl("http://myrivulet.appspot.com/");
  }
}
