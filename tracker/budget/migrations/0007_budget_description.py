# Generated by Django 5.0.1 on 2024-01-23 14:45

from django.db import migrations, models


class Migration(migrations.Migration):
    dependencies = [
        ("budget", "0006_budget_user_transaction_user"),
    ]

    operations = [
        migrations.AddField(
            model_name="budget",
            name="description",
            field=models.CharField(default="", max_length=250),
        ),
    ]
