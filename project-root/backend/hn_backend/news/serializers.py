from rest_framework import serializers
from .models import  top_news

class top_news_serializer(serializers.ModelSerializer):
  class Meta:
    model = top_news
    fields = ['Scraped_Page','Rank','Title','Link','Score','Posted_by','Time_ago','Comments']