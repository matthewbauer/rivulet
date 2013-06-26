package com.bauer.rivulet;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.view.View.OnClickListener;
import android.widget.Button;
 
public class MainActivity extends Activity {
 
	private Button button;
 
	public void onCreate(Bundle savedInstanceState) {
		final Context context = this;
 
		super.onCreate(savedInstanceState);
		setContentView(R.layout.main);
 
		button = (Button) findViewById(R.id.buttonUrl);
 
		button.setOnClickListener(new OnClickListener() {
		  @Override
		  public void onClick(View arg0) {
		    Intent intent = new Intent(context, WebViewActivity.class);
		    startActivity(intent);
		  }
		});
	}
}
