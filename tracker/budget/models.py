import calendar
import uuid
from datetime import date

from django.contrib import admin
from django.contrib.auth.models import User
from django.db import models
from django.utils import timezone


class Budget(models.Model):
    def __str__(self):
        return self.name

    unique_id = models.UUIDField(default=uuid.uuid4, editable=False, unique=True)
    name = models.CharField(max_length=150)
    created_at = models.DateTimeField(auto_now_add=True)
    user = models.ForeignKey(User, related_name="user", on_delete=models.DO_NOTHING)
    description = models.CharField(max_length=250, default="")
    shared_to = models.ManyToManyField(User, related_name="shared_to", blank=True)


class Invitation(models.Model):
    def __str__(self):
        return "%s, b(%s) | usr(%s) -> usr(%s)" % (
            self.id,
            self.budget.id,
            self.from_user.username,
            self.to_user.username,
        )

    id = models.UUIDField(
        primary_key=True, default=uuid.uuid4, editable=False, unique=True
    )
    budget = models.ForeignKey(Budget, on_delete=models.DO_NOTHING)
    from_user = models.ForeignKey(
        User, related_name="from_user", on_delete=models.DO_NOTHING
    )
    accepted = models.BooleanField(default=False)
    to_user = models.ForeignKey(
        User, related_name="to_user", on_delete=models.DO_NOTHING
    )
    token = models.UUIDField(default=uuid.uuid4, editable=False, unique=True)
    valid_from = models.DateTimeField(auto_now_add=True)
    valid_to = models.DateTimeField()

    @property
    def expired(self):
        return timezone.now() >= self.valid_to


class TransactionCategory(models.Model):
    def __str__(self):
        return self.name

    color = models.CharField(max_length=10)
    description = models.CharField(max_length=300)
    icon = models.CharField(max_length=50)
    name = models.CharField(max_length=100)
    unique_id = models.UUIDField(default=uuid.uuid4, editable=False, unique=True)


class Transaction(models.Model):
    def __str__(self):
        return "%s: %s, %s" % (self.budget.id, self.amount, self.is_expense)

    @admin.display(
        boolean=True,
        description="In Current Month?",
    )
    def in_current_month(self):
        now = date.today()
        month = calendar.monthrange(now.year, now.month)
        start_date = now.replace(day=1)
        end_date = now.replace(day=month[1])

        return start_date <= self.created_at.date() <= end_date

    budget = models.ForeignKey(Budget, on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True)
    amount = models.FloatField(null=True, blank=True, default=0.0)
    is_expense = models.BooleanField(default=True)
    description = models.CharField(max_length=300, default="")
    user = models.ForeignKey(User, on_delete=models.DO_NOTHING)
    unique_id = models.UUIDField(default=uuid.uuid4, editable=False, unique=True)
    category = models.ForeignKey(
        TransactionCategory, on_delete=models.DO_NOTHING, blank=True, null=True
    )
