/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.dask.models;

import java.util.Objects;
import java.util.Arrays;
import com.dominodatalab.distributedcomputeoperator.dask.models.V1alpha1DaskClusterSpecSchedulerConfigMapItems;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * V1alpha1DaskClusterSpecSchedulerSecret
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:36:51.907Z[Etc/UTC]")
public class V1alpha1DaskClusterSpecSchedulerSecret {
  public static final String SERIALIZED_NAME_DEFAULT_MODE = "defaultMode";
  @SerializedName(SERIALIZED_NAME_DEFAULT_MODE)
  private Integer defaultMode;

  public static final String SERIALIZED_NAME_ITEMS = "items";
  @SerializedName(SERIALIZED_NAME_ITEMS)
  private List<V1alpha1DaskClusterSpecSchedulerConfigMapItems> items = null;

  public static final String SERIALIZED_NAME_OPTIONAL = "optional";
  @SerializedName(SERIALIZED_NAME_OPTIONAL)
  private Boolean optional;

  public static final String SERIALIZED_NAME_SECRET_NAME = "secretName";
  @SerializedName(SERIALIZED_NAME_SECRET_NAME)
  private String secretName;


  public V1alpha1DaskClusterSpecSchedulerSecret defaultMode(Integer defaultMode) {
    
    this.defaultMode = defaultMode;
    return this;
  }

   /**
   * Get defaultMode
   * @return defaultMode
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Integer getDefaultMode() {
    return defaultMode;
  }


  public void setDefaultMode(Integer defaultMode) {
    this.defaultMode = defaultMode;
  }


  public V1alpha1DaskClusterSpecSchedulerSecret items(List<V1alpha1DaskClusterSpecSchedulerConfigMapItems> items) {
    
    this.items = items;
    return this;
  }

  public V1alpha1DaskClusterSpecSchedulerSecret addItemsItem(V1alpha1DaskClusterSpecSchedulerConfigMapItems itemsItem) {
    if (this.items == null) {
      this.items = new ArrayList<>();
    }
    this.items.add(itemsItem);
    return this;
  }

   /**
   * Get items
   * @return items
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public List<V1alpha1DaskClusterSpecSchedulerConfigMapItems> getItems() {
    return items;
  }


  public void setItems(List<V1alpha1DaskClusterSpecSchedulerConfigMapItems> items) {
    this.items = items;
  }


  public V1alpha1DaskClusterSpecSchedulerSecret optional(Boolean optional) {
    
    this.optional = optional;
    return this;
  }

   /**
   * Get optional
   * @return optional
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Boolean getOptional() {
    return optional;
  }


  public void setOptional(Boolean optional) {
    this.optional = optional;
  }


  public V1alpha1DaskClusterSpecSchedulerSecret secretName(String secretName) {
    
    this.secretName = secretName;
    return this;
  }

   /**
   * Get secretName
   * @return secretName
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getSecretName() {
    return secretName;
  }


  public void setSecretName(String secretName) {
    this.secretName = secretName;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1DaskClusterSpecSchedulerSecret v1alpha1DaskClusterSpecSchedulerSecret = (V1alpha1DaskClusterSpecSchedulerSecret) o;
    return Objects.equals(this.defaultMode, v1alpha1DaskClusterSpecSchedulerSecret.defaultMode) &&
        Objects.equals(this.items, v1alpha1DaskClusterSpecSchedulerSecret.items) &&
        Objects.equals(this.optional, v1alpha1DaskClusterSpecSchedulerSecret.optional) &&
        Objects.equals(this.secretName, v1alpha1DaskClusterSpecSchedulerSecret.secretName);
  }

  @Override
  public int hashCode() {
    return Objects.hash(defaultMode, items, optional, secretName);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1DaskClusterSpecSchedulerSecret {\n");
    sb.append("    defaultMode: ").append(toIndentedString(defaultMode)).append("\n");
    sb.append("    items: ").append(toIndentedString(items)).append("\n");
    sb.append("    optional: ").append(toIndentedString(optional)).append("\n");
    sb.append("    secretName: ").append(toIndentedString(secretName)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

