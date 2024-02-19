import calendar

from django.contrib import messages
from django.contrib.auth.decorators import login_required
from django.contrib.auth.mixins import LoginRequiredMixin
from django.contrib.auth.models import User
from django.db.models import Prefetch, Q
from django.shortcuts import redirect, render
from django.utils import timezone
from django.utils.decorators import method_decorator
from django.views import View, generic

from .models import Budget, Invitation, Transaction


def grouped_transaction_by_catergory(t, grouped_by: dict) -> None:
    if t["transaction__category__name"] not in grouped_by:
        grouped_by[t["transaction__category__name"]] = {
            "name": str(t["transaction__category__name"]),
            "icon": t["transaction__category__icon"],
            "color": t["transaction__category__color"],
            "description": t["transaction__category__description"],
            "id": t["transaction__category__unique_id"],
            "total": t["transaction__amount"],
        }
    else:
        amount = t["transaction__amount"]
        total = grouped_by[t["transaction__category__name"]]["total"]
        grouped_by[t["transaction__category__name"]]["total"] = round(total + amount, 2)


def calculate_transactions_flow(t, transactions_flow: dict[str, int]) -> None:
    day = t["transaction__created_at"].day
    month = t["transaction__created_at"].month
    label = f"{day}.{month}"
    transactions_flow[label] = round(
        transactions_flow[label] + t["transaction__amount"], 2
    )


def create_empty_transactions_flow(today):
    monthrange = calendar.monthrange(today.year, today.month)
    list_of_days = list(range(1, monthrange[1]))
    list_of_labels = [f"{day}.{today.month}" for day in list_of_days]
    return dict.fromkeys(list_of_labels, 0)


@method_decorator(login_required, name="dispatch")
class Dashboard(View):
    template_name = "dashboard/overall.html"

    # [TODO]: Add tests
    def get(self, request):
        today = timezone.now()
        transactions_flow = create_empty_transactions_flow(today)
        total_expenses = 0
        total_income = 0
        grouped_by = {}

        transactions = list(
            Budget.objects.filter(Q(user=request.user) | Q(shared_to=request.user))
            .prefetch_related(
                Prefetch(
                    "transaction_set",
                    queryset=Transaction.objects.filter(
                        created_at__year=str(today.year),
                        created_at__month=str(today.month),
                    ).order_by("-created_at"),
                )
            )
            .values(
                "transaction",
                "transaction__amount",
                "transaction__description",
                "transaction__created_at",
                "transaction__unique_id",
                "transaction__user",
                "transaction__user__id",
                "transaction__user__username",
                "transaction__category",
                "transaction__category__color",
                "transaction__category__description",
                "transaction__category__icon",
                "transaction__category__name",
                "transaction__category__unique_id",
            )
            .order_by("transaction__created_at")
        )

        for t in transactions:
            if t["transaction"] is not None:
                grouped_transaction_by_catergory(t, grouped_by)
                calculate_transactions_flow(t, transactions_flow)

                if t["transaction__amount"] >= 0:
                    total_income += t["transaction__amount"]
                else:
                    total_expenses += t["transaction__amount"]

        context = {
            "total": round(total_income + total_expenses, 2),
            "total_income": round(total_income, 2),
            "total_expenses": round(total_expenses, 2),
            "latest_transactions": transactions[:10],
            "transactions_per_category": grouped_by,
            "transactions_flow": transactions_flow,
        }
        return render(request, self.template_name, context=context)


@method_decorator(login_required, name="dispatch")
class BudgetList(View):
    template_name = "budget/budget-list.html"

    def get(self, request):
        """
        Budgets list for login user
        """
        shared_budgets = Budget.objects.filter(shared_to=request.user)
        invitations = Invitation.objects.filter(
            from_user=request.user, valid_to__gt=timezone.now()
        )
        budgets = Budget.objects.filter(
            Q(user=request.user) | Q(shared_to=request.user)
        )
        print(budgets)

        return render(
            request,
            self.template_name,
            {
                "budgets": budgets,
                "shared_budgets": shared_budgets,
                "invitations": invitations,
            },
        )


@method_decorator(login_required, name="dispatch")
class BudgetShareToken(View):
    template_name = "budget/budget-share.html"

    def get(self, request):
        if "invitationToken" not in request.GET:
            messages.error(request, "Token not present")
        else:
            token = request.GET["invitationToken"]
            invitation = Invitation.objects.get(token=token)

            if invitation.to_user.id == request.user.id:
                invitation.budget.shared_to.add(request.user)
                invitation.accepted = True
                invitation.save()
            else:
                messages.error(request, "Invalid user from invitation")

            return render(
                request, self.template_name, context={"invitation": invitation}
            )


@method_decorator(login_required, name="dispatch")
class BudgetDetail(View):
    template_name = "budget/budget-detail.html"

    def get(self, request, pk):
        budget = Budget.objects.get(pk=pk)
        shared_to_users = budget.shared_to.all()
        
        print(shared_to_users)
        return render(
            request,
            self.template_name,
            context={"budget": budget, "shared_to_users": shared_to_users},
        )


# [TODO] Add update function for transaction
class TransactionDetail(LoginRequiredMixin, generic.DetailView):
    model = Transaction
    template_name = "transaction/transaction-detail.html"


# class BudgetInvite(View):
#     def post(self, request, budget_id):
#         users_emails = request.GET["users_emails"]
#         emails = users_emails.split(",")
#         budget = Budget.objects.get(pk=budget_id)
#
#         for email in emails:
#             user = User.objects.get(email=email)
#             invitation = Invitation(budget=budget, to_user=user, from_user=request.user)
#             invitation.save()
#
#             msg = (
#                 """
#             User: %s invited you to budget: %s \n
#             Open link in the browser:
#
#             http://localhost:8000/join?invitationToken=%s
#             """
#                 % user.username,
#                 budget.name,
#                 invitation.token,
#             )
#
#             send_mail(
#                 subject="Invite to budget",
#                 message=msg,
#                 from_email="from@mymoney.com",
#                 recipient_list=[email],
#             )


class BudgetShare(View):
    def post(self, request, budget_id):
        # TODO: this should be 2 step feature.
        # User A types user emails in field on front, and User B should accept link in the email
        share_to = request.POST["share_to"]
        budget = Budget.objects.get(pk=budget_id)
        user_share_to = User.objects.get(pk=share_to)
        budget.shared_to.add(user_share_to)
        budget.save()
        messages.success(request, "Budget shared")
        return redirect("/budget/%s" % budget_id)


@method_decorator(login_required, name="dispatch")
class TransactionAdd(View):
    def post(self, request, budget_id):
        """
        Add new transaction to budget
        """
        amount = request.POST["amount"]
        description = request.POST["desc"]
        budget = Budget.objects.get(pk=budget_id)

        # [TODO] Add Validation
        transaction = Transaction(
            amount=amount,
            description=description,
            budget=budget,
            user=request.user,
        )
        transaction.save()
        messages.success(request, "Transaction added!")
        context = { "transaction": transaction}
        return render(request, template_name='partials/transaction.html', context=context)


@method_decorator(login_required, name="dispatch")
class BudgetAdd(View):
    template_name = "budget/budget-add.html"

    def get(self, request):
        """
        Render BudgetAdd form
        """
        return render(request, self.template_name)

    def post(self, request):
        """
        Create new budget for login user
        """
        name = request.POST["budget_name"]

        # Add proper validation
        if name == "":
            messages.error(request, "Budget name empty")
        else:
            messages.success(request, "Budget created!")
            new_budget = Budget(name=name, user=request.user)
            new_budget.save()

        return render(request, self.template_name)
