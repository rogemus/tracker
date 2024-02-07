# Generated by Django 5.0.1 on 2024-02-07 22:23

from django.db import migrations, models


class Migration(migrations.Migration):
    dependencies = [
        ("budget", "0020_remove_transactioncategory_transaction_and_more"),
    ]

    operations = [
        migrations.RemoveField(
            model_name="transaction",
            name="is_expense",
        ),
        migrations.AlterField(
            model_name="transaction",
            name="description",
            field=models.CharField(blank=True, default="", max_length=300, null=True),
        ),
    ]
