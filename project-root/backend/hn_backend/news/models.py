from django.db import models

# Create your models here.
class top_news(models.Model):
  Scraped_Page=models.URLField(max_length=300)
  Rank = models.CharField(max_length=10)
  Title = models.CharField(max_length=300)
  Link = models.URLField(max_length=300)
  Score = models.CharField(max_length=1000)
  Posted_by = models.CharField(max_length=300)
  Time_ago = models.CharField(max_length=300)
  Comments = models.CharField(max_length=300)
