from django.urls import path
from .views import add_top_news

urlpatterns = [
    # Example: path('', views.home, name='home'),
    path('top-news/',add_top_news,name='add_top_news')
]